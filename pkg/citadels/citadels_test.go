package citadels

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestTable(t *testing.T) {
	table := NewTable()

	onEv := func(logging bool) OnEventFunc {
		return func(e Event, p *Player) {
			if logging{
				t.Log(e.Type)
			}
			switch e.Type {
			case EventTypeChooseHero:
				data, ok := e.Data.(EventChooseHero)
				if !ok {
					return
				}
				table.SelectHero(p, data.Heroes[0].Name)
			case EventTypeNextTurn:
				data, ok := e.Data.(EventNextTurn)
				if !ok {
					return
				}
				if data.PlayerID == p.ID {
					table.MakeAction(ActionTypeCoin, string(p.ID))
					table.EndTurn(p.ID)
				}
			}
		}
	}

	p1 := NewPlayer(PlayerID(strconv.Itoa(rand.Intn(100000000000))), onEv(true))
	p2 := NewPlayer(PlayerID(strconv.Itoa(rand.Intn(100000000000))), onEv(false))
	p3 := NewPlayer(PlayerID(strconv.Itoa(rand.Intn(100000000000))), onEv(false))
	p4 := NewPlayer(PlayerID(strconv.Itoa(rand.Intn(100000000000))), onEv(false))

	err := table.AddPlayer(p1)
	if err != nil {
		t.Error(err)
	}
	err = table.AddPlayer(p2)
	if err != nil {
		t.Error(err)
	}
	err = table.AddPlayer(p3)
	if err != nil {
		t.Error(err)
	}
	err = table.AddPlayer(p4)
	if err != nil {
		t.Error(err)
	}
	err = table.Start()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Minute * 2)
}