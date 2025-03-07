package yanfei

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var chargeFrames []int

const chargeHitmark = 63

func init() {
	chargeFrames = frames.InitAbilSlice(79)          // CA -> N1
	chargeFrames[action.ActionCharge] = 78           // CA -> CA
	chargeFrames[action.ActionSkill] = chargeHitmark // CA -> E
	chargeFrames[action.ActionBurst] = chargeHitmark // CA -> Q
	chargeFrames[action.ActionDash] = 51             // CA -> D
	chargeFrames[action.ActionJump] = 49             // CA -> J
	chargeFrames[action.ActionSwap] = 59             // CA -> Swap
}

// Charge attack function - handles seal use
func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {

	//check for seal stacks
	if !c.StatusIsActive(sealBuffKey) {
		c.sealCount = 0
	}

	// apply a1
	c.a1(c.sealCount)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       charge[c.sealCount][c.TalentLvlAttack()],
	}

	// add windup if we're in idle or swap only
	windup := 16
	if c.Core.Player.CurrentState() == action.Idle || c.Core.Player.CurrentState() == action.SwapState {
		windup = 0
	}

	// TODO: Not sure of snapshot timing
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 2, false, combat.TargettableEnemy, combat.TargettableGadget), chargeHitmark-windup, chargeHitmark-windup)

	c.Core.Log.NewEvent("yanfei charge attack consumed seals", glog.LogCharacterEvent, c.Index).
		Write("current_seals", c.sealCount)

	// Clear the seals next frame just in case for some reason we call stam check late
	c.Core.Tasks.Add(func() {
		c.sealCount = 0
		c.DeleteStatus(sealBuffKey)
	}, 1)

	// needed for a4 hitlag handling
	c.a4HitlagApplied = false

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return chargeFrames[next] - windup },
		AnimationLength: chargeFrames[action.InvalidAction] - windup,
		CanQueueAfter:   chargeFrames[action.ActionJump] - windup, // earliest cancel is before hitmark
		State:           action.ChargeAttackState,
	}
}
