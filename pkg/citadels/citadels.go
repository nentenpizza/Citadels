// Package citadels contains internal logic of game
// exports API to interact with game world and display it
package citadels

import (
	"sync"
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
var (
	// PickPhase is phase when players selects their heroes
	PickPhase Phase = "citadels.phase.pick"
	// ActionPhase is phase when players perform actions
	ActionPhase Phase = "citadels.phase.action"
)

// Table represents a game table (also known as Room)
type Table struct {
	sync.Mutex

	// king is player which starts PickPhase
	king *Player

	// turn is player who is currently taking a turn
	turn *Player

	started bool

	currentPhase Phase

	// heroesToSelect is map of remaining heroes
	// when a Player selected a hero, the hero should disappear from the map
	// and the pick should go to the next player with the current map state
	// used only in PickPhase
	heroesToSelect map[string]Card

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

	// makes random player a king
	for _, p := range t.players {
		t.king = p
		break
	}

	return nil
}

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

