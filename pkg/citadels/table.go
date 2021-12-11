// Package citadels contains internal logic of game
// exports API to interact with game world and display it
package citadels

import (
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func init() {
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
	// EndGamePhase is phase when game is ended
	EndGamePhase Phase = "citadels.phase.end"
	PreGamePhase Phase = "citadels.phase.pregame"
)

var heroSets = map[string][]Hero{
	"default": {Witch(), Blackmailer(), Enchantress(), Emperor(), Abat(), Alchemist(), Architect(), Warlord(), CustomsOfficer()},
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

	completedQuartersFirst *Player

	// heroesToSelect is map of remaining heroes
	// when a Player selected a hero, the hero should disappear from the slice
	// and the pick should go to the next player with the current map state
	// used only in PickPhase
	heroesToSelect []Hero

	openLockedHeroes   []Hero
	closedLockedHeroes []Hero

	deck []Quarter

	bewitchedPlayer *Player

	players map[PlayerID]*Player

	Delays bool

	done chan struct{}
}

func NewTable(delay bool) *Table {
	return &Table{
		players:      make(map[PlayerID]*Player),
		currentPhase: PreGamePhase,
		Delays:       delay,
		done:         make(chan struct{}),
	}
}

func (t *Table) Ended() <-chan struct{} {
	return t.done
}

func (t *Table) close() {
	for _, p := range t.players {
		close(p.updates)
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

	var i int = 1
	for _, player := range t.players {
		player.Order = i
		i++
		go player.Listen()
	}

	// TODO: rethink in future
	t.heroSet = "default"

	// makes random player a king
	for _, p := range t.players {
		t.king = p
		break
	}

	t.doBroadcastEvent(Event{
		Type: EventTypeRevealHeroSet,
		Data: EventHeroSet{
			HeroSet: heroSets[t.heroSet],
		},
	})

	if t.Delays {
		t.Sleep(DelayAfterHeroSetReveal)
	}

	t.drawCards()

	t.started = true
	t.doBroadcastEvent(Event{
		Type: EventTypeGameStarted,
		Data: EventGameStarted{
			King: t.king,
		},
	})

	t.startPickPhase()

	return nil
}

func (t *Table) drawCards() {
	// TODO: fill the deck normally
	deck := make([]Quarter, 100)
	types := []string{
		QuarterTypeMilitary, QuarterTypeSpecial, QuarterTypeNoble, QuarterTypeSpiritual, QuarterTypeTrade,
	}
	for i := 0; i < 100; i++ {
		deck[i] = Quarter{
			Name:  strconv.Itoa(rand.Intn(10000000)),
			Type:  types[rand.Intn(5)],
			Price: 1,
			//Price: rand.Intn(5)+1,
		}
	}
	t.deck = deck
	for _, p := range t.players {
		p.AvailableQuarters = t.deck[:4]
		t.deck = t.deck[4:]
		p.Notify(Event{
			Type: EventTypeDrawCards,
			Data: EventCards{Cards: p.AvailableQuarters},
		})
	}
}

func (t *Table) startPickPhase() {
	t.currentPhase = PickPhase

	heroSet := make([]Hero, len(heroSets[t.heroSet]))
	copy(heroSet, heroSets[t.heroSet])

	rand.Shuffle(len(heroSet), func(i, j int) { heroSet[i], heroSet[j] = heroSet[j], heroSet[i] })

	t.openLockedHeroes = make([]Hero, 0)
	t.closedLockedHeroes = make([]Hero, 0)
	var finalIndex int

	for i, hero := range heroSet {
		t.openLockedHeroes = append(t.openLockedHeroes, hero)
		if len(t.players) == 4 && i+1 == 3 {
			t.closedLockedHeroes = append(t.closedLockedHeroes, heroSet[0])
			finalIndex = i
			break
		}
	}
	t.heroesToSelect = heroSet[finalIndex:]

	t.doBroadcastEvent(Event{
		Type: EventTypePickPhaseStarted,
		Data: EventPickPhaseStarted{
			OpenLockedHeroes:   t.openLockedHeroes,
			ClosedLockedHeroes: len(t.closedLockedHeroes),
		},
	})

	t.selecting = t.king
	t.king.Notify(Event{
		Type: EventTypeChooseHero,
		Data: EventChooseHero{Heroes: t.heroesToSelect},
	})

	t.startSelectingTimer()
}

func (t *Table) nextSelecting() {
	playersCount := len(t.players)
	currentPlayerOrder := t.selecting.Order

	nextPlayerOrder := currentPlayerOrder + 1
	if nextPlayerOrder > playersCount {
		nextPlayerOrder = 1
	}

	for _, p := range t.players {
		if p.Order == nextPlayerOrder {

			// if turn returns to the king this is means that all players at the table selected their heroes
			if t.king.ID == p.ID {
				t.startActionPhase()
				return
			}
			t.selecting = p

			p.Notify(Event{
				Type: EventTypeChooseHero,
				Data: EventChooseHero{Heroes: t.heroesToSelect},
			})

			t.doBroadcastEvent(Event{
				Type: EventTypeNextSelecting,
				Data: EventPlayerID{PlayerID: p.ID},
			})
			break
		}
	}

	t.startSelectingTimer()
}

func (t *Table) startActionPhase() {
	t.currentPhase = ActionPhase

	t.doBroadcastEvent(Event{
		Type: EventTypeActionPhaseStarted,
	})

	t.nextTurn()
}

func (t *Table) nextTurn() {
	t.currentIndex += 1
	if t.currentIndex > 9 {
		t.currentIndex = 1
		t.endRound()
		if t.currentPhase != EndGamePhase {
			t.startPickPhase()
		}
		return
	}

	for _, p := range t.players {
		if p.Hero.Turn == t.currentIndex {
			t.turn = p

			t.turn.madeAction = false
			t.turn.currentCardsChoice = nil
			t.turn.BuildChancesLeft = 1

			t.doBroadcastEvent(Event{
				Type: EventTypeNextTurn,
				Data: EventNextTurn{
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
		Type: EventTypeHeroIsAbsent,
		Data: EventHeroIsAbsent{
			Turn: t.currentIndex,
		},
	})

	t.Sleep(DelayAfterHeroAbsent)

	t.nextTurn()
}

func (t *Table) EndTurn(pID PlayerID) {
	t.Lock()
	defer t.Unlock()

	if t.turn.ID == pID {
		t.nextTurn()
	}
}

func (t *Table) endRound() {
	for _, p := range t.players {
		if len(p.CompletedQuarters) == 7 && t.currentPhase != EndGamePhase {
			t.currentPhase = EndGamePhase
		}
	}

	if t.currentPhase != EndGamePhase {
		return
	}
	t.completedQuartersFirst.totalScore += 4
	var winner *Player
	for _, p := range t.players {
		for _, quarter := range p.CompletedQuarters {
			p.totalScore += quarter.Price
		}
		if len(p.CompletedQuarters) == 7 && p.ID != t.completedQuartersFirst.ID {
			p.totalScore += 2
		}
		if winner == nil {
			winner = p
		}
		if p.totalScore > winner.totalScore {
			winner = p
		}
	}

	if winner == nil {
		return
	}

	t.doBroadcastEvent(Event{Type: EventTypeGameEnded,
		Data: EventGameEnded{Winner: winner.ID},
	})

	t.close()
	t.done <- struct{}{}
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
func (t *Table) King() *Player {
	t.Lock()
	defer t.Unlock()
	return t.king
}

// Turn returns player which is currently taking a turn
func (t *Table) Turn() *Player {
	t.Lock()
	defer t.Unlock()
	return t.turn
}

// SelectHero sets Player.Hero to Hero with given heroName if
// Player is currently selecting and Hero with heroName is present in heroesToSelect
func (t *Table) SelectHero(p *Player, heroName string) {
	t.Lock()
	defer t.Unlock()

	if t.currentPhase != PickPhase {
		return
	}

	if t.selecting.ID != p.ID {
		p.Notify(Event{Error: ErrorTypeAnotherPlayerSelecting})
		return
	}

	for i, hero := range t.heroesToSelect {
		if hero.Name == heroName {
			t.selecting.Hero = t.heroesToSelect[i]
			t.heroesToSelect = removeHero(t.heroesToSelect, i)
			p.Notify(Event{
				Type: EventTypeHeroSelected,
				Data: EventHeroSelected{Hero: hero},
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

func (t *Table) CastSkill(casterID string, ev Event) error {
	t.Lock()
	defer t.Unlock()

	caster, ok := t.PlayerByID(casterID)
	if !ok {
		return ErrPlayerNotExists
	}

	err := caster.Hero.Skill.Do(t, caster, ev)
	if err != nil {
		return err
	}

	t.nextTurn()
	return nil
}

const (
	ActionTypeCoin  = "coins"
	ActionTypeCards = "cards"
)

// MakeAction gives Player 2 coins or 1 card depending on the type
func (t *Table) MakeAction(actionType string, pID string) {
	t.Lock()
	defer t.Unlock()

	target, ok := t.playerByID(pID)
	if !ok {
		return
	}

	if t.turn.ID != target.ID {
		return
	}

	if target.madeAction {
		return
	}

	switch actionType {
	case ActionTypeCoin:
		target.AddCoins(2)
		t.doBroadcastEvent(Event{Type: EventTypeCoinsGive, Data: EventCoinGive{
			To:     target.ID,
			Amount: 2,
			Sum:    target.Coins,
		}})

	case ActionTypeCards:
		target.setCurrentCardsChoice(t.deck[:2])
		target.Notify(Event{Type: EventTypeChooseCards, Data: EventChooseCards{
			Cards: t.deck[:2],
		}})
		t.deck = t.deck[2:]
		t.doBroadcastEvent(Event{Type: EventTypePlayerChoosingCards, Data: EventPlayerChoosingCards{
			PlayerID:    target.ID,
			CardsAmount: 2,
		}})
	default:
		target.Notify(Event{
			Error: ErrorTypeWrongAction,
		})
		return
	}

	target.madeAction = true
}

// SelectCard adds card to Player.AvailableQuarters
func (t *Table) SelectCard(cardName string, pID string) {
	t.Lock()
	defer t.Unlock()
	target, ok := t.playerByID(pID)
	if !ok {
		return
	}

	if target.currentCardsChoice == nil || len(target.currentCardsChoice) < 2 {
		return
	}

	for i, card := range target.currentCardsChoice {
		if card.Name == cardName {
			target.AddQuarter(card)
			target.currentCardsChoice = nil

			t.doBroadcastEvent(Event{Type: EventTypePlayerSelectedCard, Data: EventPlayerSelectedCard{
				PlayerID: target.ID,
				Index:    i,
			}})
			break
		}
	}

}

func (t *Table) BuildQuarter(quarter Quarter, pID string) {
	t.Lock()
	defer t.Unlock()

	target, ok := t.playerByID(pID)
	if !ok {
		return
	}

	if target.BuildChancesLeft < 1 {
		return
	}

	if t.currentPhase != ActionPhase {
		return
	}

	if t.turn.ID != target.ID {
		return
	}

	if quarter.Price > target.Coins {
		target.Notify(Event{
			Error: ErrorTypeNotEnoughCoins,
		})
		return
	}

	if !target.hasQuarter(quarter.Name) {
		return
	}

	if target.builtQuarter(quarter.Name) {
		target.Notify(Event{
			Error: ErrorTypeQuarterAlreadyBuilt,
		})
		return
	}

	target.buildQuarter(quarter)
	t.doBroadcastEvent(Event{Type: EventTypePlayerBuiltQuarter,
		Data: EventQuarter{Quarter: quarter},
	})
	target.SubtractBuildChancesLeft(1)

	if len(target.CompletedQuarters) == 7 {
		t.completedQuartersFirst = target
	}
}

// AddPlayer adds player to the table
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

// RemovePlayer removes player from the table
func (t *Table) RemovePlayer(pID PlayerID) error {
	t.Lock()
	defer t.Unlock()
	if t.started {
		return ErrTableAlreadyStarted
	}
	p := t.players[pID]
	delete(t.players, pID)
	p.Table = nil
	return nil
}

func (t *Table) Sleep(dur time.Duration) {
	if t.Delays {
		time.Sleep(dur)
	}
}
