package citadels

import "errors"

var (
	ErrTableIsFull         = errors.New("table is full")
	ErrNotEnoughPlayers    = errors.New("not enough players")
	ErrTableAlreadyStarted = errors.New("table already started")

	ErrWrongEventData = errors.New("wrong event data")
	ErrPlayerNotExists = errors.New("player does not exists")
	ErrCannotCastOnMyself = errors.New("can not cast on myself")
)

// Errors for events
var (
	ErrorTypeTargetHasNoCoins = "errors.target.no.coins"
	ErrorTypeTargetHasNoCards = "errors.target.no.cards"
	ErrorTypeAnotherPlayerSelecting = "errors.selecting.another"
	ErrorTypeHeroNotInStack = "errors.hero.not.stack"
)