// Package citadels contains internal logic of game
// exports API to interact with game world and display it
package citadels

import (
	"math/rand"
	"sync"
	"time"
)

func init(){
	rand.Seed(time.Now().UnixNano())
}
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

	currentIndex int

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

	// TODO: rethink in future
	t.heroSet = "default"

	// makes random player a king
	for _, p := range t.players {
		t.king = p
		break
	}

	t.doBroadcastEvent(Event{
		Type:  EventTypeRevealHeroSet,
		Data:  EventHeroSet{
			HeroSet: heroSets[t.heroSet],
		},
	})

	time.Sleep(DelayAfterHeroSetReveal)

	t.started = true
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

	heroSet := make([]Hero, 0)
	copy(heroSet, heroSets[t.heroSet])
	rand.Shuffle(len(heroSet), func(i, j int) { heroSet[i], heroSet[j] = heroSet[j], heroSet[i] })

	t.openLockedHeroes = make([]Hero, 0)
	t.closedLockedHeroes = make([]Hero, 0)
	var finalIndex int

	for i, hero := range heroSet{
		t.openLockedHeroes = append(t.openLockedHeroes, hero)
		if len(t.players) == 4 && i+1 == 3 {
			t.closedLockedHeroes = append(t.closedLockedHeroes, heroSet[0])
			finalIndex = i
			break
		}
	}

	t.heroesToSelect = heroSet[finalIndex:]

	t.doBroadcastEvent(Event{
		Type:  EventTypePickPhaseStarted,
		Data: EventPickPhaseStarted{
			OpenLockedHeroes:   t.openLockedHeroes,
			ClosedLockedHeroes: len(t.closedLockedHeroes),
		},
	})

	t.selecting = t.king
	t.king.Notify(Event{
		Type:  EventTypeChooseHero,
		Data:  EventChooseHero{ Heroes: t.heroesToSelect },
	})

	t.startSelectingTimer()
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

			// if turn returns to the king this is means that all players at the table selected their heroes
			if t.king.ID == p.ID{
				t.startActionPhase()
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

	t.startSelectingTimer()
}

func (t *Table) startActionPhase() {
	t.currentPhase = ActionPhase

	t.doBroadcastEvent(Event{
		Type:  EventTypeActionPhaseStarted,
	})

	t.nextTurn()

	t.startTurnTimer()
}


func (t *Table) nextTurn() {
	t.currentIndex += 1
	if t.currentIndex > 9 {
		t.currentIndex = 1
		t.startPickPhase()
		return
	}

	for _, p := range t.players{
		if p.Hero.Turn == t.currentIndex{
			t.turn = p
			t.doBroadcastEvent(Event{
				Type:  EventTypeNextTurn,
				Data:  EventNextTurn{
					PlayerID: p.ID,
					Hero:     p.Hero,
					Turn:     p.Hero.Turn,
				},
			})

			t.startTurnTimer()
			return
		}
	}

	t.doBroadcastEvent(Event{
		Type:  EventTypeHeroIsAbsent,
		Data:  EventHeroIsAbsent{
			Turn: t.currentIndex,
		},
	})

	time.Sleep(DelayAfterHeroAbsent)

	t.nextTurn()
}

func (t *Table) startTurnTimer() {
	playerID := t.turn.ID
	time.AfterFunc(time.Minute, func() {
		t.Lock()
		defer t.Unlock()
		if t.currentPhase != ActionPhase {
			return
		}
		if t.turn.ID == playerID {
			t.nextTurn()
		}
	})
}

func (t *Table) startSelectingTimer() {
	playerID := t.selecting.ID
	time.AfterFunc(time.Minute, func() {
		t.Lock()
		defer t.Unlock()
		if t.currentPhase != ActionPhase {
			return
		}
		if t.selecting.ID == playerID {
			t.forceSelecting()
		}
	})
}

// Started returns current state of the table
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

// Turn returns player which is currently taking a turn
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
	p.Table = t
	return nil
}

