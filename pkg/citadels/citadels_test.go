package citadels

import (
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestTable(t *testing.T) {
	table := NewTable()

	onEvLog := func(e Event, p *Player){
		log.Println(e.Type)
		switch e.Type {
		case EventTypeChooseHero:
			data, ok := e.Data.(EventChooseHero)
			if !ok{
				return
			}
			table.SelectHero(p, data.Heroes[0].Name)
		case EventTypeNextTurn:
			data, ok := e.Data.(EventNextTurn)
			if !ok{
				return
			}
			if data.PlayerID == p.ID{
				table.MakeAction(ActionTypeCoin, string(p.ID))
				table.EndTurn(p.ID)
			}
		}
	}

	onEv := func(e Event, p *Player){
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

	p1 := NewPlayer(PlayerID(strconv.Itoa(rand.Intn(100000000000))), onEvLog)
	p2 := NewPlayer(PlayerID(strconv.Itoa(rand.Intn(100000000000))), onEv)
	p3 := NewPlayer(PlayerID(strconv.Itoa(rand.Intn(100000000000))), onEv)
	p4 := NewPlayer(PlayerID(strconv.Itoa(rand.Intn(100000000000))), onEv)

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