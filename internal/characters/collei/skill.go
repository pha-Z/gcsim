package collei

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

const (
	skillKey     = "collei-skill"
	skillRelease = 20
	skillReturn  = 157
)

var (
	skillHitmarks = []int{34, 138}
	skillFrames   []int
)

func init() {
	skillFrames = frames.InitAbilSlice(68)
	skillFrames[action.ActionAttack] = 65
	skillFrames[action.ActionAim] = 65
	skillFrames[action.ActionSkill] = 67
	skillFrames[action.ActionDash] = 54
	skillFrames[action.ActionJump] = 53
	skillFrames[action.ActionSwap] = 66
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	// The game has ICD as AttackTagElementalArt, ICDTagElementalArt,
	// ICDGroupColleiBoomerangForward, and ICDGroupColleiBoomerangBack. However,
	// we believe this is unnecessary, so just use ICDTagNone.
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Floral Brush",
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeSlash,
		Element:            attributes.Dendro,
		Durability:         25,
		Mult:               skill[c.TalentLvlSkill()],
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}
	var c6Cb func(a combat.AttackCB)
	if c.Base.Cons >= 6 {
		c6Triggered := false
		c6Cb = func(_ combat.AttackCB) {
			if c6Triggered {
				return
			}
			c6Triggered = true
			c.c6()
		}
	}
	for _, hitmark := range skillHitmarks {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 2, false, combat.TargettableEnemy),
			skillRelease,
			hitmark,
			c6Cb,
		)
	}

	c.Core.Tasks.Add(func() {
		c.AddStatus(skillKey, skillReturn-skillRelease, false)
	}, skillRelease)

	c.sproutShouldExtend = false
	c.sproutShouldProc = c.Base.Cons >= 2
	c.Core.Tasks.Add(func() {
		if !c.sproutShouldProc {
			return
		}
		src := c.Core.F
		c.sproutSrc = src
		duration := 180
		if c.sproutShouldExtend {
			duration += 180
		}
		c.AddStatus(sproutKey, duration, true)
		ai := c.a1AttackInfo()
		snap := c.Snapshot(&ai)
		c.QueueCharTask(func() {
			c.a1Ticks(src, snap)
		}, sproutHitmark)
	}, skillReturn)

	c.Core.QueueParticle("collei", 3, attributes.Dendro, skillHitmarks[0]+c.ParticleDelay)

	c.SetCDWithDelay(action.ActionSkill, 720, 20)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionJump], // earliest cancel
		State:           action.SkillState,
	}
}
