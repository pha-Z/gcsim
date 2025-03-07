package tighnari

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var aimedFrames []int
var aimedWreathFrames []int

const aimedHitmark = 86
const aimedWreathHitmark = 175

func init() {
	aimedFrames = frames.InitAbilSlice(94)
	aimedFrames[action.ActionDash] = aimedHitmark
	aimedFrames[action.ActionJump] = aimedHitmark

	aimedWreathFrames = frames.InitAbilSlice(183)
	aimedWreathFrames[action.ActionDash] = aimedWreathHitmark
	aimedWreathFrames[action.ActionJump] = aimedWreathHitmark
}

func (c *char) Aimed(p map[string]int) action.ActionInfo {
	level, ok := p["level"]
	if !ok {
		level = 0
	}

	if c.StatusIsActive(vijnanasuffusionStatus) {
		level = 1
	}
	if level == 1 {
		return c.WreathAimed(p)
	}

	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot := p["weakspot"]

	ai := combat.AttackInfo{
		ActorIndex:           c.Index,
		Abil:                 "Aim (Charged)",
		AttackTag:            combat.AttackTagExtra,
		ICDTag:               combat.ICDTagNone,
		ICDGroup:             combat.ICDGroupDefault,
		Element:              attributes.Dendro,
		Durability:           25,
		Mult:                 fullaim[c.TalentLvlAttack()],
		HitWeakPoint:         weakspot == 1,
		HitlagHaltFrames:     0.12 * 60,
		HitlagFactor:         0.01,
		HitlagOnHeadshotOnly: true,
		IsDeployable:         true,
	}

	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), .1, false, combat.TargettableEnemy), aimedHitmark, aimedHitmark+travel)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(aimedFrames),
		AnimationLength: aimedFrames[action.InvalidAction],
		CanQueueAfter:   aimedHitmark,
		State:           action.AimState,
	}
}

func (c *char) WreathAimed(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	wreathTravel, ok := p["wreath"]
	if !ok {
		wreathTravel = 35
	}
	weakspot := p["weakspot"]

	skip := 0
	if c.StatusIsActive(vijnanasuffusionStatus) {
		skip = 142 // 2.4 * 60

		arrows := c.Tag(wreatharrows) - 1
		c.SetTag(wreatharrows, arrows)
		if arrows == 0 {
			c.DeleteStatus(vijnanasuffusionStatus)
		}
	}
	if c.Base.Cons >= 6 {
		skip += 0.9 * 60
	}
	if skip > aimedWreathHitmark {
		skip = aimedWreathHitmark
	}

	ai := combat.AttackInfo{
		ActorIndex:           c.Index,
		Abil:                 "Wreath Arrow",
		AttackTag:            combat.AttackTagExtra,
		ICDTag:               combat.ICDTagNone,
		ICDGroup:             combat.ICDGroupDefault,
		Element:              attributes.Dendro,
		Durability:           25,
		Mult:                 wreath[c.TalentLvlAttack()],
		HitWeakPoint:         weakspot == 1,
		HitlagHaltFrames:     0.12 * 60,
		HitlagFactor:         0.01,
		HitlagOnHeadshotOnly: true,
		IsDeployable:         true,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), .1, false, combat.TargettableEnemy), aimedWreathHitmark-skip, aimedWreathHitmark+travel-skip)
	c.Core.Tasks.Add(c.a1, aimedWreathHitmark-skip+1)

	ai = combat.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Clusterbloom Arrow",
		AttackTag:    combat.AttackTagExtra,
		ICDTag:       combat.ICDTagExtraAttack,
		ICDGroup:     combat.ICDGroupTighnari,
		Element:      attributes.Dendro,
		Durability:   25,
		Mult:         clusterbloom[c.TalentLvlAttack()],
		HitWeakPoint: false, // TODO: tighnari can hit the weak spot on some enemies (like hilichurls)
	}
	c.Core.Tasks.Add(func() {
		snap := c.Snapshot(&ai)
		for i := 0; i < 4; i++ {
			c.Core.QueueAttackWithSnap(
				ai,
				snap,
				combat.NewCircleHit(c.Core.Combat.Player(), 0.1, false, combat.TargettableEnemy),
				wreathTravel,
			)
		}

		if c.Base.Cons >= 6 {
			ai = combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Karma Adjudged From the Leaden Fruit",
				AttackTag:  combat.AttackTagExtra,
				ICDTag:     combat.ICDTagNone,
				ICDGroup:   combat.ICDGroupDefault,
				Element:    attributes.Dendro,
				Durability: 25,
				Mult:       1.5,
			}
			c.Core.QueueAttackWithSnap(
				ai,
				snap,
				combat.NewCircleHit(c.Core.Combat.Player(), 0.1, false, combat.TargettableEnemy),
				wreathTravel,
			)
		}
	}, aimedWreathHitmark+travel-skip)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return aimedWreathFrames[next] - skip },
		AnimationLength: aimedWreathFrames[action.InvalidAction] - skip,
		CanQueueAfter:   aimedWreathHitmark - skip,
		State:           action.AimState,
	}
}
