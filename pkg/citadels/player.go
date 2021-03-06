package citadels

import (
	"log"
	"sync"
)

type OnEventFunc func(Event, *Player)

type PlayerID string

type Player struct {
	sync.Mutex
	// Quarters that player have and can build
	AvailableQuarters []Quarter `json:"-"`

	// Quarters that player already built
	CompletedQuarters []Quarter `json:"completed_quarters"`

	// BuildChancesLeft shows how many Quarter`s player can build this round
	BuildChancesLeft int

	// Hero that the player chooses each round
	// Has unique spells which can change outcome of the game
	Hero Hero `json:"-"`

	// In-game currency for building quarters
	Coins int `json:"coins"`

	// Order represents the way players sits at the table
	Order int

	ID PlayerID `json:"id"`

	Table *Table `json:"-"`

	madeAction bool

	currentCardsChoice []Quarter

	totalScore int

	updates chan Event

	OnEvent OnEventFunc
}

func NewPlayer(id PlayerID, onEvent OnEventFunc) *Player {
	return &Player{
		ID:      id,
		updates: make(chan Event, 10000000),
		currentCardsChoice: make([]Quarter, 0),
		AvailableQuarters: make([]Quarter, 0),
		CompletedQuarters: make([]Quarter, 0),
		OnEvent: onEvent,
	}
}

func (p *Player) Listen() {
	for {
		select {
		case e, ok := <- p.updates:
			if !ok{
				log.Println("leave listen ", p.ID)
				return
			}
			p.OnEvent(e, p)
		}
	}
}

func (p *Player) TotalScore() int  {
	return p.totalScore
}

func (p *Player) builtQuarter(quarterName string) bool {
	p.Lock()
	defer p.Unlock()
	for _, quarter := range p.CompletedQuarters{
		if quarter.Name == quarterName{
			return true
		}
	}
	return false
}

func (p *Player) hasQuarter(quarterName string) bool {
	p.Lock()
	defer p.Unlock()
	for _, quarter := range p.AvailableQuarters{
		if quarter.Name == quarterName{
			return true
		}
	}
	return false
}

func (p *Player) buildQuarter(quarter Quarter) {
	p.Lock()
	defer p.Unlock()
	p.CompletedQuarters = append(p.CompletedQuarters, quarter)
	p.AvailableQuarters = removeQuarterByName(p.AvailableQuarters, quarter.Name)
	p.BuildChancesLeft--
}

func (p *Player) SubtractBuildChancesLeft(i int) {
	p.Lock()
	defer p.Unlock()
	p.BuildChancesLeft =- i
}

func (p *Player) Updates() <-chan Event {
	return p.updates
}

func (p *Player) Notify(e Event){
	p.Lock()
	defer p.Unlock()
	p.updates <- e
}

func (p *Player) AddCoins(coins int){
	p.Lock()
	defer p.Unlock()
	p.Coins += coins
}

func (p *Player) AddQuarter(quarter Quarter){
	p.Lock()
	defer p.Unlock()
	p.AvailableQuarters = append(p.AvailableQuarters, quarter)
}

func (p *Player) setCurrentCardsChoice(cards []Quarter) {
	p.Lock()
	defer p.Unlock()
	p.currentCardsChoice = cards
}

func (p *Player) giveCoins(other *Player,coins int) {
	p.Lock()
	defer p.Unlock()

	if p.Coins <= 0 || p.Coins - coins <= 0{
		return
	}

	p.Coins -= coins
	other.Coins += coins

	ev := Event{Type: EventTypeStealCoinPrivate, Data: EventSteal{
		FromID: p.ID,
		To:     other.ID,
		Count:  coins,
	}}

	p.updates <- ev
	other.Notify(ev)

	p.Table.BroadcastEvent(Event{
		Type: EventTypeStealCoin,
		Data: EventSteal{To: other.ID, FromID: p.ID, Count: coins},
	})
}

func (p *Player) giveRandomCards(other *Player, cards int) {
	p.Lock()
	defer p.Unlock()

	if p.ID == other.ID{
		return
	}
	if len(p.AvailableQuarters) <= 0 {
		return
	}

	for i:=0;i<cards;i++{
		if len(p.AvailableQuarters) >= i+1 {
			other.Lock()
			other.AvailableQuarters = append(other.AvailableQuarters, p.AvailableQuarters[i])
			p.AvailableQuarters = removeQuarter(p.AvailableQuarters, i)
			other.Unlock()
		}
	}

	p.updates <- Event{Type: EventTypeStealCardPrivate, Data: EventStealCards{
		FromID: p.ID,
		To:     other.ID,
		AvailableQuarters:  p.AvailableQuarters,
	}}

	other.Notify(Event{Type: EventTypeStealCardPrivate, Data: EventStealCards{
		FromID: p.ID,
		To:     other.ID,
		AvailableQuarters:  other.AvailableQuarters,
	}})

	p.Table.BroadcastEvent(Event{
		Type: EventTypeStealCard,
		Data: EventSteal{To: other.ID, FromID: p.ID, Count: cards},
	})
}
