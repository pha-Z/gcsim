package damage

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/reactable"
	"github.com/genshinsim/gcsim/pkg/stats"
)

var reactions = map[event.Event]combat.ReactionType{
	event.OnOverload:           combat.Overload,
	event.OnSuperconduct:       combat.Superconduct,
	event.OnMelt:               combat.Melt,
	event.OnVaporize:           combat.Vaporize,
	event.OnFrozen:             combat.Freeze,
	event.OnElectroCharged:     combat.ElectroCharged,
	event.OnSwirlHydro:         combat.SwirlHydro,
	event.OnSwirlCryo:          combat.SwirlCryo,
	event.OnSwirlElectro:       combat.SwirlElectro,
	event.OnSwirlPyro:          combat.SwirlPyro,
	event.OnCrystallizeCryo:    combat.CrystallizeCryo,
	event.OnCrystallizeElectro: combat.CrystallizeElectro,
	event.OnCrystallizeHydro:   combat.CrystallizeHydro,
	event.OnCrystallizePyro:    combat.CrystallizePyro,
	event.OnAggravate:          combat.Aggravate,
	event.OnSpread:             combat.Spread,
	event.OnQuicken:            combat.Quicken,
	event.OnBloom:              combat.Bloom,
	event.OnHyperbloom:         combat.Hyperbloom,
	event.OnBurgeon:            combat.Burgeon,
	event.OnBurning:            combat.Burning,
}

func init() {
	stats.Register(NewStat)
}

type buffer struct {
	damageOverTime        map[string]float64
	damageByChar          []map[string]float64
	damageByCharByTargets []map[string]float64
	damageInstancesByChar []map[string]int
	abilUsageCountByChar  []map[string]int
	charActiveTime        []int
	elementUptime         []map[string]int
	particleCount         map[string]float64
	reactionsTriggered    map[string]int
}

func NewStat(core *core.Core) (stats.StatsCollector, error) {
	out := buffer{
		damageOverTime:        make(map[string]float64),
		damageByChar:          make([]map[string]float64, len(core.Player.Chars())),
		damageByCharByTargets: make([]map[string]float64, len(core.Player.Chars())),
		damageInstancesByChar: make([]map[string]int, len(core.Player.Chars())),
		abilUsageCountByChar:  make([]map[string]int, len(core.Player.Chars())),
		charActiveTime:        make([]int, len(core.Player.Chars())),
		elementUptime:         make([]map[string]int, len(core.Combat.Enemies())),
		particleCount:         make(map[string]float64),
		reactionsTriggered:    make(map[string]int),
	}
	var sb strings.Builder

	for i := 0; i < len(core.Player.Chars()); i++ {
		out.damageByChar[i] = make(map[string]float64)
		out.damageByCharByTargets[i] = make(map[string]float64)
		out.damageInstancesByChar[i] = make(map[string]int)
		out.abilUsageCountByChar[i] = make(map[string]int)
	}

	for i := 0; i < len(core.Combat.Enemies()); i++ {
		out.elementUptime[i] = make(map[string]int)
	}

	core.Events.Subscribe(event.OnActionExec, func(args ...interface{}) bool {
		active := args[0].(int)
		action := args[1].(action.Action)
		out.abilUsageCountByChar[active][action.String()]++
		return false
	}, "legacy-sim-abil-usage")

	core.Events.Subscribe(event.OnParticleReceived, func(args ...interface{}) bool {
		p := args[0].(character.Particle)
		out.particleCount[p.Source] += p.Num
		return false
	}, "legacy-particles-log")

	eventSubFunc := func(reaction combat.ReactionType) func(args ...interface{}) bool {
		return func(args ...interface{}) bool {
			out.reactionsTriggered[string(reaction)]++
			return false
		}
	}

	for k, v := range reactions {
		core.Events.Subscribe(k, eventSubFunc(v), "legacy-reaction-log")
	}

	core.Events.Subscribe(event.OnTick, func(args ...interface{}) bool {
		out.charActiveTime[core.Player.Active()]++
		for i, e := range core.Combat.Enemies() {
			if t, ok := e.(*enemy.Enemy); ok {
				for r, v := range t.Durability {
					if v > reactable.ZeroDur {
						out.elementUptime[i][reactable.ReactableModifier(r).String()]++
					}
				}
			}
		}
		return false
	}, "legacy-on-tick")

	core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		t := args[0].(combat.Target)
		atk := args[1].(*combat.AttackEvent)
		dmg := args[2].(float64)

		// No need to pull damage stats for non-enemies
		if t.Type() != combat.TargettableEnemy {
			return false
		}

		if atk.Info.DoNotLog {
			return false
		}

		sb.Reset()
		sb.WriteString(atk.Info.Abil)
		if atk.Info.Amped {
			if atk.Info.AmpMult == 1.5 {
				sb.WriteString(" [amp: 1.5]")
			} else if atk.Info.AmpMult == 2 {
				sb.WriteString(" [amp: 2.0]")
			}
		}

		if atk.Info.Catalyzed {
			if atk.Info.CatalyzedType == combat.Aggravate {
				sb.WriteString(" [aggravate]")
			} else if atk.Info.CatalyzedType == combat.Spread {
				sb.WriteString(" [spread]")
			}
		}

		out.damageByChar[atk.Info.ActorIndex][sb.String()] += dmg
		out.damageByCharByTargets[atk.Info.ActorIndex][strconv.Itoa(t.Index())] += dmg
		if dmg > 0 {
			out.damageInstancesByChar[atk.Info.ActorIndex][sb.String()] += 1
		}

		frameBucket := fmt.Sprintf("%.2f", float64(int(core.F/15)*15)/60.0)
		out.damageOverTime[frameBucket] += dmg

		return false
	}, "legacy-dmg-log")

	return &out, nil
}

func (b buffer) Flush(core *core.Core, result *stats.Result) {
	result.Legacy.DamageOverTime = b.damageOverTime
	result.Legacy.DamageByChar = b.damageByChar
	result.Legacy.DamageByCharByTargets = b.damageByCharByTargets
	result.Legacy.DamageInstancesByChar = b.damageInstancesByChar
	result.Legacy.AbilUsageCountByChar = b.abilUsageCountByChar
	result.Legacy.CharActiveTime = b.charActiveTime
	result.Legacy.ElementUptime = b.elementUptime
	result.Legacy.ParticleCount = b.particleCount
	result.Legacy.ReactionsTriggered = b.reactionsTriggered
}
