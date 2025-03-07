package frostbearer

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

func init() {
	core.RegisterWeaponFunc(keys.Frostbearer, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	atk := 0.65 + float64(r)*0.15
	atkc := 1.6 + float64(r)*0.4
	prob := 0.5 + float64(r)*0.1

	const icdKey = "frostbearer-icd"
	icd := 600

	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		if ae.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		if ae.Info.AttackTag != combat.AttackTagNormal && ae.Info.AttackTag != combat.AttackTagExtra {
			return false
		}
		if c.Rand.Float64() < prob {
			char.AddStatus(icdKey, icd, true)

			ai := combat.AttackInfo{
				ActorIndex: char.Index,
				Abil:       "Frostbearer Proc",
				AttackTag:  combat.AttackTagWeaponSkill,
				ICDTag:     combat.ICDTagNone,
				ICDGroup:   combat.ICDGroupDefault,
				Element:    attributes.Physical,
				Durability: 100,
				Mult:       atk,
			}

			if t.AuraContains(attributes.Cryo) || t.AuraContains(attributes.Frozen) {
				ai.Mult = atkc
			}

			c.QueueAttack(ai, combat.NewCircleHit(t, 3, false, combat.TargettableEnemy), 0, 1)

		}
		return false
	}, fmt.Sprintf("forstbearer-%v", char.Base.Key.String()))

	return w, nil
}
