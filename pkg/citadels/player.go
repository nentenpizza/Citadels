package citadels

type PlayerID string

type Player struct {
	// Quarters that player have and can build
	HandQuarters []Card

	// Quarters that player already built
	CompletedQuarters []Card

	// Hero that the player chooses each round
	// Has unique spells which can change outcome of the game
	Hero Card

	// In-game currency for building quarters
	Coins int

	ID PlayerID

	updates chan *Event
}

func NewPlayer(id PlayerID) *Player {
	return &Player{
		ID:      id,
		updates: make(chan *Event),
	}
}

func (p *Player) Updates() <-chan *Event {
	return p.updates
}
