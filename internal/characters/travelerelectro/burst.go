package travelerelectro

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var burstFrames [][]int

const burstHitmark = 37

func init() {
	burstFrames = make([][]int, 2)

	// Male
	burstFrames[0] = frames.InitAbilSlice(63) // Q -> E
	burstFrames[0][action.ActionAttack] = 62  // Q -> N1
	burstFrames[0][action.ActionDash] = 62    // Q -> D
	burstFrames[0][action.ActionJump] = 61    // Q -> J
	burstFrames[0][action.ActionSwap] = 60    // Q -> Swap

	// Female
	burstFrames[1] = frames.InitAbilSlice(62) // Q -> E/D
	burstFrames[1][action.ActionAttack] = 61  // Q -> N1
	burstFrames[1][action.ActionJump] = 61    // Q -> J
	burstFrames[1][action.ActionSwap] = 61    // Q -> Swap
}

/*
*
[12:01 PM] pai: never tried to measure it but emc burst looks like it has roughly 1~1.5 abyss tile of range, skill goes a bit further i think
[12:01 PM] pai: the 3 hits from the skill also like split out and kind of auto target if that's useful information
*
*/
func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Bellowing Thunder",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 5, false, combat.TargettableEnemy, combat.TargettableGadget), 0, burstHitmark)

	c.SetCDWithDelay(action.ActionBurst, 1200, 35)
	c.ConsumeEnergy(37)

	// emc burst is not hitlag extendable
	c.Core.Status.Add("travelerelectroburst", 720) // 12s, starts on cast

	procAI := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Falling Thunder Proc (Q)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       burstTick[c.TalentLvlBurst()],
	}
	c.burstSnap = c.Snapshot(&procAI)
	c.burstAtk = &combat.AttackEvent{
		Info:     procAI,
		Snapshot: c.burstSnap,
	}
	c.burstSrc = c.Core.F

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames[c.gender]),
		AnimationLength: burstFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   burstFrames[c.gender][action.ActionJump], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) burstProc() {
	icd := 0

	// Lightning Shroud
	//  When your active character's Normal or Charged Attacks hit opponents, they will call Falling Thunder forth, dealing Electro DMG.
	//  When Falling Thunder hits opponents, it will regenerate Energy for that character.
	//  One instance of Falling Thunder can be generated every 0.5s.
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		t := args[0].(combat.Target)

		// only apply on na/ca
		if ae.Info.AttackTag != combat.AttackTagNormal && ae.Info.AttackTag != combat.AttackTagExtra {
			return false
		}
		// make sure the person triggering the attack is on field still
		if ae.Info.ActorIndex != c.Core.Player.Active() {
			return false
		}
		// only apply if burst is active
		if c.Core.Status.Duration("travelerelectroburst") == 0 {
			return false
		}
		// One instance of Falling Thunder can be generated every 0.5s.
		if icd > c.Core.F {
			c.Core.Log.NewEvent("travelerelectro Q (active) on icd", glog.LogCharacterEvent, c.Index)
			return false
		}

		// Use burst snapshot, update target & source frame
		atk := *c.burstAtk
		atk.SourceFrame = c.Core.F
		//attack is 2 (or 2.5 for enhanced) aoe centered on target
		atk.Pattern = combat.NewCircleHit(t, 2, false, combat.TargettableEnemy, combat.TargettableGadget)

		// C2 - Violet Vehemence
		// When Falling Thunder created by Bellowing Thunder hits an opponent, it will decrease their Electro RES by 15% for 8s.
		// c6 - World-Shaker
		//  Every 2 Falling Thunder attacks triggered by Bellowing Thunder will greatly increase the DMG
		//  dealt by the next Falling Thunder, which will deal 200% of its original DMG and will restore
		//  an additional 1 Energy to the current character.
		c.c6Damage(&atk)
		atk.Callbacks = append(atk.Callbacks, c.fallingThunderEnergy(), c.c2(t), c.c6Energy())

		c.Core.QueueAttackEvent(&atk, 1)

		c.Core.Log.NewEvent("travelerelectro Q proc'd", glog.LogCharacterEvent, c.Index).
			Write("char", ae.Info.ActorIndex).
			Write("attack tag", ae.Info.AttackTag)

		icd = c.Core.F + 30 // 0.5s
		return false
	}, "travelerelectro-bellowingthunder")
}

func (c *char) fallingThunderEnergy() combat.AttackCBFunc {
	return func(_ combat.AttackCB) {
		// Regenerate 1 flat energy for the active character
		active := c.Core.Player.ActiveChar()
		active.AddEnergy("travelerelectro-fallingthunder", burstRegen[c.TalentLvlBurst()])
	}
}
