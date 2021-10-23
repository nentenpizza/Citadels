package citadels

import "errors"

var (
	ErrTableIsFull         = errors.New("table is full")
	ErrNotEnoughPlayers    = errors.New("not enough players")
	ErrTableAlreadyStarted = errors.New("table already started")
)
