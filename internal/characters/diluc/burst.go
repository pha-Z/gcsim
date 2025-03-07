package diluc

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const burstHitmark = 100

func init() {
	burstFrames = frames.InitAbilSlice(140) // Q -> D
	burstFrames[action.ActionAttack] = 139  // Q -> N1
	burstFrames[action.ActionSkill] = 139   // Q -> E
	burstFrames[action.ActionJump] = 139    // Q -> J
	burstFrames[action.ActionSwap] = 138    // Q -> Swap
}

const burstBuffKey = "diluc-q"

func (c *char) phoenixDMG(ai combat.AttackInfo, dot int, explode int) func() {
	return func() {
		// DoT does damage every .2 seconds for 7 hits? so every 12 frames
		// DoT does max 7 hits + explosion, roughly every 13 frame? blows up at 210 frames
		// DoT
		for i := 0; i < dot; i++ {
			c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 2, false, combat.TargettableEnemy, combat.TargettableGadget), 0, i*12)
		}
		// Explosion
		if explode > 0 {
			ai.Abil = "Dawn (Explode)"
			ai.Mult = burstExplode[c.TalentLvlBurst()]
			c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 2, false, combat.TargettableEnemy, combat.TargettableGadget), 0, 98)
		}
	}
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	dot, ok := p["dot"]
	if !ok {
		dot = 2 //number of dot hits
	}
	if dot > 7 {
		dot = 7
	}
	explode, ok := p["explode"]
	if !ok {
		explode = 0 //if explode hits
	}

	//enhance weapon for 12 seconds (with a4)
	// Infusion starts when burst starts and ends when burst comes off CD - check any diluc video
	c.AddStatus(burstBuffKey, 720, true)

	// a4: add 20% pyro damage
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(burstBuffKey, 720),
		AffectedStat: attributes.PyroP,
		Amount: func() ([]float64, bool) {
			return c.a4buff, true
		},
	})

	// Snapshot occurs late in the animation when it is released from the claymore
	// For our purposes, snapshot upon damage proc
	c.Core.Tasks.Add(func() {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               "Dawn (Strike)",
			AttackTag:          combat.AttackTagElementalBurst,
			ICDTag:             combat.ICDTagElementalBurst,
			ICDGroup:           combat.ICDGroupDiluc,
			StrikeType:         combat.StrikeTypeBlunt,
			Element:            attributes.Pyro,
			Durability:         50,
			Mult:               burstInitial[c.TalentLvlBurst()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   0.09 * 60,
			CanBeDefenseHalted: true,
		}

		c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 2, false, combat.TargettableEnemy, combat.TargettableGadget), 0, 1)

		// both initial hit, DoT and explosion all have 50 durability
		ai.Abil = "Dawn (Tick)"
		ai.Mult = burstDOT[c.TalentLvlBurst()]

		// only initial hit has hitlag
		ai.HitlagHaltFrames = 0
		ai.CanBeDefenseHalted = false

		// TODO: also consider making this actually sort of move (like aoe wise)
		// queue DoT and Explosion DMG
		c.QueueCharTask(c.phoenixDMG(ai, dot, explode), 12)
	}, burstHitmark)

	c.ConsumeEnergy(21)
	c.SetCDWithDelay(action.ActionBurst, 720, 14)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}
