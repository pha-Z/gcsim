package crescent

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterWeaponFunc(keys.CrescentPike, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	atk := .15 + float64(r)*.05
	const buffKey = "crescent-pike-buff"
	buffDuration := 300 // 5s * 60

	c.Events.Subscribe(event.OnParticleReceived, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}
		c.Log.NewEvent("crescent pike active", glog.LogWeaponEvent, char.Index).
			Write("expiry (without hitlag)", c.F+300)
		char.AddStatus(buffKey, buffDuration, true)

		return false
	}, fmt.Sprintf("cp-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		if ae.Info.ActorIndex != char.Index {
			return false
		}
		if ae.Info.AttackTag != combat.AttackTagNormal && ae.Info.AttackTag != combat.AttackTagExtra {
			return false
		}
		if char.StatusIsActive(buffKey) {
			ai := combat.AttackInfo{
				ActorIndex: char.Index,
				Abil:       "Crescent Pike Proc",
				AttackTag:  combat.AttackTagWeaponSkill,
				ICDTag:     combat.ICDTagNone,
				ICDGroup:   combat.ICDGroupDefault,
				Element:    attributes.Physical,
				Durability: 100,
				Mult:       atk,
			}
			trg := args[0].(combat.Target)
			c.QueueAttack(ai, combat.NewCircleHit(trg, 0.1, false, combat.TargettableEnemy), 0, 1)
		}
		return false
	}, fmt.Sprintf("cpp-%v", char.Base.Key.String()))
	return w, nil
}
