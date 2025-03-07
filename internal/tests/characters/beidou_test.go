package characters

import (
	"testing"

	_ "github.com/genshinsim/gcsim/internal/characters/beidou"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

// Test to make sure in 2 target scenario, Beidou burst bounces between the 2 targets
func TestBeidouBounce(t *testing.T) {
	c, trg := makeCore(2)
	prof := defProfile(keys.Beidou)
	prof.Base.Cons = 6
	idx, err := c.AddChar(prof)
	if err != nil {
		t.Errorf("error adding char: %v", err)
		t.FailNow()
	}
	c.Player.SetActive(idx)
	err = c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}
	//initialize some settings
	c.Combat.DefaultTarget = trg[0].Key()
	c.QueueParticle("system", 1000, attributes.NoElement, 0)
	advanceCoreFrame(c)

	//start tests
	dmgCount := make(map[combat.TargetKey]int)
	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		ae, ok := args[1].(*combat.AttackEvent)
		if !ok {
			return false
		}
		if ae.Info.Abil == "Stormbreak Proc (Q)" {
			dmgCount[t.Key()]++
		}

		return false
	}, "q-bounce-count")

	p := make(map[string]int)
	c.Player.Exec(action.ActionBurst, keys.Beidou, p)
	for !c.Player.CanQueueNextAction() {
		advanceCoreFrame(c)
	}
	done := false
	for !done {
		err := c.Player.Exec(action.ActionAttack, keys.Beidou, p)
		switch err {
		case player.ErrActionNotReady, player.ErrPlayerNotReady, player.ErrActionNoOp:
			advanceCoreFrame(c)
		case nil:
			done = true
		default:
			t.Errorf("encountered unexpected error: %v", err)
			t.FailNow()
		}
	}
	for i := 0; i < 200; i++ {
		advanceCoreFrame(c)
	}

	if dmgCount[trg[0].Key()] != 3 {
		t.Errorf("expecting target 0 (key %v) to take 3 hits, got %v", trg[0].Key(), dmgCount[trg[0].Key()])
	}

	if dmgCount[trg[1].Key()] != 2 {
		t.Errorf("expecting target 1 (key %v) to take 2 hits, got %v", trg[1].Key(), dmgCount[trg[1].Key()])
	}

}
