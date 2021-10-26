// Package citadels contains internal logic of game
// exports API to interact with game world and display it
package citadels

import (
	"sync"
	"time"
)

const (
	// MaxPlayers is max number of players in 1 room
	MaxPlayers = 4
	// MinPlayers is min number of players in 1 room
	MinPlayers = 4
)

// Phase represents separate logic cycles in game
// for example, in PickPhase players should pick a hero and move on to the next phase
type Phase string

// Phases
const (
	// PickPhase is phase when players selects their heroes
	PickPhase Phase = "citadels.phase.pick"
	// ActionPhase is phase when players perform actions
	ActionPhase Phase = "citadels.phase.action"
)

var heroSets = map[string][]Hero{
	"default": {Emperor(), Emperor(), Emperor(), Emperor(), Emperor(), Emperor()},
}

// Table represents a game table (also known as Room)
type Table struct {
	sync.Mutex

	// king is player which starts PickPhase
	king *Player

	// turn is player who is currently taking a turn
	turn *Player

	// selecting is Player who selecting Hero right now
	selecting *Player

	started bool

	currentPhase Phase

	heroSet string

	// heroesToSelect is map of remaining heroes
	// when a Player selected a hero, the hero should disappear from the slice
	// and the pick should go to the next player with the current map state
	// used only in PickPhase
	heroesToSelect []Hero

	openLockedHeroes []Hero
	closedLockedHeroes []Hero

	players map[PlayerID]*Player
}

func NewTable() *Table {
	return &Table{
		players: make(map[PlayerID]*Player),
	}
}

// Start makes the table ready to conduct rounds
func (t *Table) Start() error {
	t.Lock()
	defer t.Unlock()
	if t.started {
		return ErrTableAlreadyStarted
	}
	if len(t.players) < MinPlayers {
		return ErrNotEnoughPlayers
	}
	if len(t.players) > MaxPlayers {
		return ErrTableIsFull
	}
	t.started = true

	// TODO: rethink in future
	t.heroSet = "default"

	// makes random player a king
	for _, p := range t.players {
		t.king = p
		break
	}

	t.doBroadcastEvent(Event{
		Type:  EventTypeGameStarted,
		Data:  EventGameStarted{
			King: t.king,
		},
	})

	t.startPickPhase()

	return nil
}

func (t *Table) startPickPhase()  {
	t.currentPhase = PickPhase
	t.selecting = t.king

	copy(t.heroesToSelect, heroSets[t.heroSet])

	t.doBroadcastEvent(Event{
		Type:  EventTypePickPhaseStarted,
	})

	t.king.Notify(Event{
		Type:  EventTypeChooseHero,
		Data:  EventChooseHero{ Heroes: t.heroesToSelect },
	})

	selectingID := t.selecting.ID
	time.AfterFunc(time.Minute, func(){
		t.Lock()
		defer t.Unlock()
		if t.selecting.ID == selectingID{
			t.forceSelecting()
		}
	})
}

func (t *Table) nextSelecting(){
	playersCount := len(t.players)
	currentPlayerOrder := t.selecting.Order

	nextPlayerOrder := currentPlayerOrder + 1
	if nextPlayerOrder > playersCount {
		nextPlayerOrder = 1
	}

	for _, p := range t.players{
		if p.Order == nextPlayerOrder{

			if t.king.ID == p.ID{
				t.endPickPhase()
				return
			}
			t.selecting = p

			p.Notify(Event{
				Type:  EventTypeChooseHero,
				Data:  EventChooseHero{ Heroes: t.heroesToSelect },
			})

			t.doBroadcastEvent(Event{
				Type:  EventTypeNextSelecting,
				Data:  EventPlayerID{PlayerID: p.ID},
			})
			break
		}
	}

	selectingID := t.selecting.ID
	time.AfterFunc(time.Minute, func(){
		t.Lock()
		defer t.Unlock()
		if t.selecting.ID == selectingID{
			t.forceSelecting()
		}
	})

}

func (t *Table) endPickPhase() {
	t.currentPhase = ActionPhase

	t.doBroadcastEvent(Event{
		Type:  EventTypeActionPhaseStarted,
	})

}

// Started returns current state of table
func (t *Table) Started() bool {
	t.Lock()
	defer t.Unlock()
	return t.started
}

// King returns player which starts PickPhase this round
func (t *Table) King() *Player  {
	t.Lock()
	defer t.Unlock()
	return t.king
}

// Turn returns player who is currently taking a turn
func (t *Table) Turn() *Player  {
	t.Lock()
	defer t.Unlock()
	return t.turn
}

// SelectHero sets Player.Hero to Hero with given heroName if
// Player is currently selecting and Hero with heroName is present in heroesToSelect
func (t *Table) SelectHero(p *Player, heroName string){
	t.Lock()
	defer t.Unlock()

	if t.selecting.ID != p.ID{
		p.Notify(Event{Error: ErrorTypeAnotherPlayerSelecting})
	}

	for i, hero := range t.heroesToSelect{
		if hero.Name == heroName{
			t.selecting.Hero = t.heroesToSelect[i]
			t.heroesToSelect = removeHero(t.heroesToSelect, i)
			p.Notify(Event{
				Type:  EventTypeHeroSelected,
				Data:  EventHeroSelected{Hero: hero},
			})
			t.nextSelecting()
			return
		}
	}
	p.Notify(Event{
		Data:  EventChooseHero{Heroes: t.heroesToSelect},
		Error: ErrorTypeHeroNotInStack,
	})
}

// AddPlayer adds player to table, returns nil if success
func (t *Table) AddPlayer(p *Player) error {
	t.Lock()
	defer t.Unlock()
	if t.started {
		return ErrTableAlreadyStarted
	}
	if len(t.players) == MaxPlayers {
		return ErrTableIsFull
	}
	t.players[p.ID] = p
	return nil
}

