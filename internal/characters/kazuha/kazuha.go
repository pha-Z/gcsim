package kazuha

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

func init() {
	core.RegisterCharFunc(keys.Kazuha, NewChar)
}

type char struct {
	*tmpl.Character
	a1Absorb              attributes.Element
	a1AbsorbCheckLocation combat.AttackPattern
	qAbsorb               attributes.Element
	qFieldSrc             int
	qAbsorbCheckLocation  combat.AttackPattern
	c2buff                []float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = normalHitNum

	c.a1AbsorbCheckLocation = combat.NewCircleHit(c.Core.Combat.Player(), 1.5, false, combat.TargettableEnemy, combat.TargettablePlayer, combat.TargettableGadget)
	c.qAbsorbCheckLocation = combat.NewCircleHit(c.Core.Combat.Player(), 1.5, false, combat.TargettableEnemy, combat.TargettablePlayer, combat.TargettableGadget)

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()

	// make sure to use the same key everywhere so that these passives don't stack
	c.Core.Player.AddStamPercentMod("utility-dash", -1, func(a action.Action) (float64, bool) {
		if a == action.ActionDash && c.HPCurrent > 0 {
			return -0.2, false
		}
		return 0, false
	})

	if c.Base.Cons >= 2 {
		c.c2buff = make([]float64, attributes.EndStatType)
		c.c2buff[attributes.EM] = 200
	}
	return nil
}
