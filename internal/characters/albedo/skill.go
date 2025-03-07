package albedo

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var skillFrames []int

const skillHitmark = 25

func init() {
	skillFrames = frames.InitAbilSlice(33) // E -> Q
	skillFrames[action.ActionAttack] = 32  // E -> N1
	skillFrames[action.ActionDash] = 29    // E -> D
	skillFrames[action.ActionJump] = 28    // E -> J
	skillFrames[action.ActionSwap] = 31    // E -> Swap
}

const (
	skillICDKey = "albedo-skill-icd"
)

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Abiogenesis: Solar Isotoma",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	//TODO: damage frame
	c.bloomSnapshot = c.Snapshot(&ai)
	c.Core.QueueAttackWithSnap(ai, c.bloomSnapshot, combat.NewCircleHit(c.Core.Combat.Player(), 3, false, combat.TargettableEnemy), skillHitmark)

	//snapshot for ticks
	ai.Abil = "Abiogenesis: Solar Isotoma (Tick)"
	ai.ICDTag = combat.ICDTagElementalArt
	ai.Mult = skillTick[c.TalentLvlSkill()]
	ai.UseDef = true
	c.skillAttackInfo = ai
	c.skillSnapshot = c.Snapshot(&c.skillAttackInfo)

	//create a construct
	// Construct is not fully formed until after the hit lands (exact timing unknown)
	c.Core.Tasks.Add(func() {
		c.Core.Constructs.New(c.newConstruct(1800), true)
		c.lastConstruct = c.Core.F
		c.skillActive = true
		// Reset ICD after construct is created
		c.DeleteStatus(skillICDKey)
		// add C4 and C6 checks
		if c.Base.Cons >= 4 {
			c.Core.Tasks.Add(c.c4(c.Core.F), 18) // start checking in 0.3s
		}
		if c.Base.Cons >= 6 {
			c.Core.Tasks.Add(c.c6(c.Core.F), 18) // start checking in 0.3s
		}
	}, skillHitmark)

	c.SetCDWithDelay(action.ActionSkill, 240, 23)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillHitmark,
		State:           action.SkillState,
	}
}

func (c *char) skillHook() {
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		if !c.skillActive {
			return false
		}
		if c.StatusIsActive(skillICDKey) {
			return false
		}
		// Can't be triggered by itself when refreshing
		if atk.Info.Abil == "Abiogenesis: Solar Isotoma" {
			return false
		}

		// this ICD is most likely tied to the construct, so it's not hitlag extendable
		c.AddStatus(skillICDKey, 120, false) // proc every 2s

		snap := c.skillSnapshot

		// a1: skill tick deal 25% more dmg if enemy hp < 50%
		if c.Core.Combat.DamageMode && t.HP()/t.MaxHP() < .5 {
			snap.Stats[attributes.DmgP] += 0.25
			c.Core.Log.NewEvent("a1 proc'd, dealing extra dmg", glog.LogCharacterEvent, c.Index).
				Write("hp %", t.HP()/t.MaxHP()).
				Write("final dmg", snap.Stats[attributes.DmgP])
		}

		c.Core.QueueAttackWithSnap(c.skillAttackInfo, snap, combat.NewCircleHit(c.Core.Combat.Player(), 3, false, combat.TargettableEnemy), 1)

		//67% chance to generate 1 geo orb
		if c.Core.Rand.Float64() < 0.67 {
			c.Core.QueueParticle("albedo", 1, attributes.Geo, c.ParticleDelay)
		}

		// c1: skill tick regen 1.2 energy
		if c.Base.Cons >= 1 {
			c.AddEnergy("albedo-c1", 1.2)
			c.Core.Log.NewEvent("c1 restoring energy", glog.LogCharacterEvent, c.Index)
		}

		// c2: skill tick grant stacks, lasts 30s; each stack increase burst dmg by 30% of def, stack up to 4 times
		if c.Base.Cons >= 2 {
			if !c.StatusIsActive(c2key) {
				c.c2stacks = 0
			}
			c.AddStatus(c2key, 1800, true) //lasts 30 sec
			c.c2stacks++
			if c.c2stacks > 4 {
				c.c2stacks = 4
			}
		}

		return false
	}, "albedo-skill")
}
