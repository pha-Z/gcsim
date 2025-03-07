package yoimiya

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var burstFrames []int

const burstHitmark = 75
const abDebuff = "aurous-blaze"
const abIcdKey = "aurous-blaze-icd"

func init() {
	burstFrames = frames.InitAbilSlice(113) // Q -> N1
	burstFrames[action.ActionSkill] = 112   // Q -> E
	burstFrames[action.ActionDash] = 111    // Q -> D
	burstFrames[action.ActionJump] = 112    // Q -> J
	burstFrames[action.ActionSwap] = 109    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	//assume it does skill dmg at end of it's animation
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Aurous Blaze",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Pyro,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}
	// A4
	c.Core.Tasks.Add(c.a4, burstHitmark)

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 5, false, combat.TargettableEnemy, combat.TargettableGadget),
		0,
		burstHitmark,
		c.applyAB, // callback to apply Aurous Blaze
	)

	//add cooldown to sim
	c.SetCD(action.ActionBurst, 15*60)
	//use up energy
	c.ConsumeEnergy(5)

	c.abApplied = false

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) applyAB(a combat.AttackCB) {
	// marker an opponent after first hit
	// ignore the bouncing around for now (just assume it's always target 0)
	// icd of 2s, removed if down

	// do nothing if ab already applied on enemy
	if c.abApplied {
		return
	}
	c.abApplied = true

	trg, ok := a.Target.(*enemy.Enemy)
	// do nothing if not an enemy
	if !ok {
		return
	}

	duration := 600
	if c.Base.Cons >= 1 {
		duration = 840
	}
	trg.AddStatus(abDebuff, duration, true) // apply Aurous Blaze
}

func (c *char) burstHook() {
	//check on attack landed for target 0
	//if aurous active then trigger dmg if not on cd
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		trg, ok := args[0].(*enemy.Enemy)
		// ignore if not an enemy
		if !ok {
			return false
		}
		// ignore if debuff not on enemy
		if !trg.StatusIsActive(abDebuff) {
			return false
		}
		// ignore for self
		if ae.Info.ActorIndex == c.Index {
			return false
		}
		//ignore if on icd
		if trg.StatusIsActive(abIcdKey) {
			return false
		}
		//ignore if wrong tags
		switch ae.Info.AttackTag {
		case combat.AttackTagNormal:
		case combat.AttackTagExtra:
		case combat.AttackTagPlunge:
		case combat.AttackTagElementalArt:
		case combat.AttackTagElementalBurst:
		default:
			return false
		}
		//do explosion, set icd
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Aurous Blaze (Explode)",
			AttackTag:  combat.AttackTagElementalBurst,
			ICDTag:     combat.ICDTagElementalBurst,
			ICDGroup:   combat.ICDGroupDefault,
			Element:    attributes.Pyro,
			Durability: 25,
			Mult:       burstExplode[c.TalentLvlBurst()],
		}
		c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 3, false, combat.TargettableEnemy, combat.TargettableGadget), 0, 1)

		trg.AddStatus(abIcdKey, 120, true) // trigger Aurous Blaze ICD

		// C4
		if c.Base.Cons >= 4 {
			c.ReduceActionCooldown(action.ActionSkill, 72)
		}

		return false

	}, "yoimiya-burst-check")

	if c.Core.Flags.DamageMode {
		//add check for if yoimiya dies
		c.Core.Events.Subscribe(event.OnCharacterHurt, func(_ ...interface{}) bool {
			if c.HPCurrent <= 0 {
				// remove Aurous Blaze from target
				for _, x := range c.Core.Combat.Enemies() {
					trg := x.(*enemy.Enemy)
					if trg.StatusIsActive(abDebuff) {
						trg.DeleteStatus(abDebuff)
					}
				}
			}
			return false
		}, "yoimiya-died")
	}
}
