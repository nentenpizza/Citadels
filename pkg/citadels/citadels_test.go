package citadels

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

func TestTable(t *testing.T) {
	table := NewTable(false)
	//done := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(4)
	onEv := func(logging bool) OnEventFunc {
		return func(e Event, p *Player) {
			if logging{
				t.Log(e.Type)
			}
			switch e.Type {
			case EventTypeChooseHero:
				data, ok := e.Data.(EventChooseHero)
				if !ok {
					t.Fatal("wrong event")
				}
				table.SelectHero(p, data.Heroes[0].Name)
			case EventTypeNextTurn:
				data, ok := e.Data.(EventNextTurn)
				if !ok {
					t.Fatal("wrong event")
				}
				if data.PlayerID == p.ID {
					if len(p.AvailableQuarters) < 1 && p.Coins > 0 {
						table.MakeAction(ActionTypeCards, string(p.ID))
					} else {
						table.MakeAction(ActionTypeCoin, string(p.ID))
					}
					if p.Coins > 0 && len(p.AvailableQuarters) > 0 {
						table.BuildQuarter(p.AvailableQuarters[0], string(p.ID))
						table.EndTurn(p.ID)
					}
				}
			case EventTypeChooseCards:
				data, ok := e.Data.(EventChooseCards)
				if !ok {
					t.Fatal("wrong event")
				}
				table.SelectCard(data.Cards[0].Name, string(p.ID))
				if p.Coins > 0 && len(p.AvailableQuarters) > 0 {
					table.BuildQuarter(p.AvailableQuarters[0], string(p.ID))
					if logging {
						t.Log(len(p.AvailableQuarters))
					}
				}
				table.EndTurn(p.ID)
			case EventTypeGameEnded:
				data, ok := e.Data.(EventGameEnded)
					if !ok {
						t.Fatal("wrong event")
					}
					if p.ID == data.Winner{
						t.Log("winner " + data.Winner, "total score ", p.TotalScore())
					}
				t.Log("player " + p.ID, "total score ", p.TotalScore())
				wg.Done()
			}
		}
	}

	p1 := NewPlayer(PlayerID(strconv.Itoa(rand.Intn(100000000000))), onEv(false))
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

	wg.Wait()
}