package simulation

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func SetupTargetsInCore(core *core.Core, p core.Coord, targets []enemy.EnemyProfile) error {

	// s.stats.ElementUptime = make([]map[core.EleType]int, len(s.C.Targets))
	// s.stats.ElementUptime[0] = make(map[core.EleType]int)

	if p.R == 0 {
		return errors.New("player cannot have 0 radius")
	}
	player := avatar.New(core, p.X, p.Y, p.R)
	core.Combat.SetPlayer(player)

	// add targets
	for i, v := range targets {
		if v.Pos.R == 0 {
			return fmt.Errorf("target cannot have 0 radius (index %v): %v", i, v)
		}
		e := enemy.New(core, v)
		core.Combat.AddEnemy(e)
		//s.stats.ElementUptime[i+1] = make(map[core.EleType]int)
	}

	//default target is closest to player?
	trgs := core.Combat.EnemyByDistance(p.X, p.Y, combat.InvalidTargetKey)
	core.Combat.DefaultTarget = core.Combat.Enemy(trgs[0]).Key()

	return nil
}

func SetupCharactersInCore(core *core.Core, chars []profile.CharacterProfile, initial keys.Char) error {
	if len(chars) > 4 {
		return errors.New("cannot have more than 4 characters per team")
	}
	dup := make(map[keys.Char]bool)

	active := -1
	for _, v := range chars {
		i, err := core.AddChar(v)
		if err != nil {
			return err
		}

		if v.Base.Key == initial {
			core.Player.SetActive(i)
			active = i
		}

		if _, ok := dup[v.Base.Key]; ok {
			return fmt.Errorf("duplicated character %v", v.Base.Key)
		}
		dup[v.Base.Key] = true
	}

	if active == -1 {
		return errors.New("no active character set")
	}

	return nil
}

func SetupResonance(s *core.Core) {
	chars := s.Player.Chars()
	if len(chars) < 4 {
		return //no resonance if less than 4 chars
	}
	//count number of ele first
	count := make(map[attributes.Element]int)
	for _, c := range chars {
		count[c.Base.Element]++
	}

	for k, v := range count {
		if v >= 2 {
			switch k {
			case attributes.Pyro:
				val := make([]float64, attributes.EndStatType)
				val[attributes.ATKP] = 0.25
				f := func() ([]float64, bool) {
					return val, true
				}
				for _, c := range chars {
					c.AddStatMod(character.StatMod{
						Base:         modifier.NewBase("pyro-res", -1),
						AffectedStat: attributes.NoStat,
						Amount:       f,
					})
				}
			case attributes.Hydro:
				//TODO: reduce pyro duration not implemented; may affect bennett Q?
				val := make([]float64, attributes.EndStatType)
				val[attributes.HPP] = 0.25
				for _, c := range chars {
					c.AddStatMod(character.StatMod{
						Base:         modifier.NewBase("hydro-res-hpp", -1),
						AffectedStat: attributes.HPP,
						Amount: func() ([]float64, bool) {
							return val, true
						},
					})
				}
			case attributes.Cryo:
				val := make([]float64, attributes.EndStatType)
				val[attributes.CR] = .15
				f := func(ae *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					r, ok := t.(*enemy.Enemy)
					if !ok {
						return nil, false
					}
					if r.AuraContains(attributes.Cryo) || r.AuraContains(attributes.Frozen) {
						return val, true
					}
					return nil, false
				}
				for _, c := range chars {
					c.AddAttackMod(character.AttackMod{
						Base:   modifier.NewBase("cryo-res", -1),
						Amount: f,
					})
				}
			case attributes.Electro:
				last := 0
				recover := func(args ...interface{}) bool {
					if s.F-last < 300 && last != 0 { // every 5 seconds
						return false
					}
					s.Player.DistributeParticle(character.Particle{
						Source: "electro-res",
						Num:    1,
						Ele:    attributes.Electro,
					})
					last = s.F
					return false
				}
				s.Events.Subscribe(event.OnOverload, recover, "electro-res")
				s.Events.Subscribe(event.OnSuperconduct, recover, "electro-res")
				s.Events.Subscribe(event.OnElectroCharged, recover, "electro-res")
				s.Events.Subscribe(event.OnQuicken, recover, "electro-res")
				s.Events.Subscribe(event.OnAggravate, recover, "electro-res")
				s.Events.Subscribe(event.OnHyperbloom, recover, "electro-res")
			case attributes.Geo:
				//Increases shield strength by 15%. Additionally, characters protected by a shield will have the
				//following special characteristics:

				//	DMG dealt increased by 15%, dealing DMG to enemies will decrease their Geo RES by 20% for 15s.
				f := func() (float64, bool) { return 0.15, true }
				s.Player.Shields.AddShieldBonusMod("geo-res", -1, f)

				//shred geo res of target
				s.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
					t, ok := args[0].(*enemy.Enemy)
					if !ok {
						return false
					}
					atk := args[1].(*combat.AttackEvent)
					if s.Player.Shields.PlayerIsShielded() && s.Player.Active() == atk.Info.ActorIndex {
						t.AddResistMod(enemy.ResistMod{
							Base:  modifier.NewBase("geo-res", 15*60),
							Ele:   attributes.Geo,
							Value: -0.2,
						})
					}
					return false
				}, "geo res")

				val := make([]float64, attributes.EndStatType)
				val[attributes.DmgP] = .15
				atkf := func(ae *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					if s.Player.Shields.PlayerIsShielded() && s.Player.Active() == ae.Info.ActorIndex {
						return val, true
					}
					return nil, false
				}
				for _, c := range chars {
					c.AddAttackMod(character.AttackMod{
						Base:   modifier.NewBase("geo-res", -1),
						Amount: atkf,
					})
				}

			case attributes.Anemo:
				s.Player.AddStamPercentMod("anemo-res-stam", -1, func(a action.Action) (float64, bool) {
					return -0.15, false
				})
				// TODO: movement spd increase?
				for _, c := range chars {
					c.AddCooldownMod(character.CooldownMod{
						Base:   modifier.NewBase("anemo-res-cd", -1),
						Amount: func(a action.Action) float64 { return -0.05 },
					})
				}
			case attributes.Dendro:
				val := make([]float64, attributes.EndStatType)
				val[attributes.EM] = 50
				for _, c := range chars {
					c.AddStatMod(character.StatMod{
						Base:         modifier.NewBase("dendro-res-50", -1),
						AffectedStat: attributes.EM,
						Amount: func() ([]float64, bool) {
							return val, true
						},
					})
				}

				twoBuff := make([]float64, attributes.EndStatType)
				twoBuff[attributes.EM] = 30
				twoEl := func(args ...interface{}) bool {
					for _, c := range chars {
						c.AddStatMod(character.StatMod{
							Base:         modifier.NewBaseWithHitlag("dendro-res-30", 6*60),
							AffectedStat: attributes.EM,
							Amount: func() ([]float64, bool) {
								return twoBuff, true
							},
						})
					}
					return false
				}
				s.Events.Subscribe(event.OnBurning, twoEl, "dendro-res")
				s.Events.Subscribe(event.OnBloom, twoEl, "dendro-res")
				s.Events.Subscribe(event.OnQuicken, twoEl, "dendro-res")

				threeBuff := make([]float64, attributes.EndStatType)
				threeBuff[attributes.EM] = 20
				threeEl := func(args ...interface{}) bool {
					for _, c := range chars {
						c.AddStatMod(character.StatMod{
							Base:         modifier.NewBaseWithHitlag("dendro-res-20", 6*60),
							AffectedStat: attributes.EM,
							Amount: func() ([]float64, bool) {
								return threeBuff, true
							},
						})
					}
					return false
				}
				s.Events.Subscribe(event.OnAggravate, threeEl, "dendro-res")
				s.Events.Subscribe(event.OnSpread, threeEl, "dendro-res")
				s.Events.Subscribe(event.OnHyperbloom, threeEl, "dendro-res")
				s.Events.Subscribe(event.OnBurgeon, threeEl, "dendro-res")
			}
		}
	}
}

func SetupMisc(c *core.Core) {
	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		//dmg tag is superconduct, target is enemy
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.AttackTag != combat.AttackTagSuperconductDamage {
			return false
		}
		//add shred
		t.AddResistMod(enemy.ResistMod{
			Base:  modifier.NewBaseWithHitlag("superconduct-phys-shred", 12*60),
			Ele:   attributes.Physical,
			Value: -0.4,
		})
		return false
	}, "superconduct")
}

func (s *Simulation) handleEnergy() {
	//energy once interval=300 amount=1 #once at frame 300
	if s.cfg.Energy.Active && s.cfg.Energy.Once {
		f := s.cfg.Energy.Start
		s.cfg.Energy.Active = false
		s.C.Tasks.Add(func() {
			s.C.Player.DistributeParticle(character.Particle{
				Source: "enemy",
				Num:    float64(s.cfg.Energy.Amount),
				Ele:    attributes.NoElement,
			})
		}, f)
		s.C.Log.NewEventBuildMsg(glog.LogEnergyEvent, -1, "energy queued (once)").
			Write("last", s.cfg.Energy.LastEnergyDrop).
			Write("cfg", s.cfg.Energy).
			Write("amt", s.cfg.Energy.Amount).
			Write("energy_frame", s.C.F+f)
	}
	//energy every interval=300,600 amount=1 #randomly every 300 to 600 frames
	if s.cfg.Energy.Active && s.C.F-s.cfg.Energy.LastEnergyDrop >= s.cfg.Energy.Start {
		f := s.C.Rand.Intn(s.cfg.Energy.End - s.cfg.Energy.Start)
		s.cfg.Energy.LastEnergyDrop = s.C.F + f
		s.C.Tasks.Add(func() {
			s.C.Player.DistributeParticle(character.Particle{
				Source: "drop",
				Num:    float64(s.cfg.Energy.Amount),
				Ele:    attributes.NoElement,
			})
		}, f)
		s.C.Log.NewEventBuildMsg(glog.LogEnergyEvent, -1, "energy queued").
			Write("last", s.cfg.Energy.LastEnergyDrop).
			Write("cfg", s.cfg.Energy).
			Write("amt", s.cfg.Energy.Amount).
			Write("energy_frame", s.C.F+f)
	}
}
