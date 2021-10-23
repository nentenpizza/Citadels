// package citadels contains internal logic of game
// exports API to interact with game world and display it
package citadels

import (
	"sync"
)

const (
	// max number of players in 1 room
	MaxPlayers = 4
	// min number of players in 1 room
	MinPlayers = 4
)

// Phases
var (
	// Phase when players selects their heroes
	PickPhase = "citadels.phase.pick"
	// Phase when players perform actions
	ActionPhase = "citadels.phase.action"
)

// Table represents a game table (also known as Room)
type Table struct {
	sync.Mutex

	// King is player which starts PickPhase
	King *Player

	// Turn is player who is currently taking a turn
	Turn *Player

	started bool

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
		t.King = p
		break
	}

	return nil
}

func (t *Table) Started() bool {
	return t.started
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
