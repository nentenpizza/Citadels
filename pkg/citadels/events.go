package citadels

// Event Types
var (
	EventTypeGameStarted = "GameStarted"
	EventTypePickPhaseStarted = "PhasePickStarted"
	EventTypeActionPhaseStarted = "ActionPhaseStarted"

	EventTypeNextSelecting = "NextSelecting"

	EventTypeSelectHero = "SelectHero"

	EventTypeCastSkill = "CastSkill"

	EventTypeStealCoin = "StealCoin"
	EventTypeStealCard = "StealCard"
	EventTypeStealCoinPrivate = "StealCoinPrivate"
	EventTypeStealCardPrivate = "StealCardPrivate"

	EventTypeHeroSelected = "hero.selected"
	EventTypeChooseHero = "ChooseHero"
)

type Event struct {
	Type string `json:"type"`
	Data interface{} `json:"data"`
	Error string `json:"error"`
}

// Events Data
type(
	EventTargeted struct {
		TargetID string `json:"target_id"`
	}

	EventEmperorSkill struct {
		TargetID PlayerID `json:"target_id"`

		// If Coin is false it means that player wants card instead of coin
		Coin bool `json:"coin"`
	}

	EventSteal struct {
		// FromID is who gives coin/card away
		FromID PlayerID `json:"from_id"`

		// To is receiver
		To PlayerID `json:"to"`
		Count int `json:"count"`
	}

	EventStealCards struct {
		// FromID is who gives cards away
		FromID PlayerID `json:"from_id"`

		// To is receiver
		To PlayerID `json:"to"`

		// New info about  Player.AvailableQuarters
		AvailableQuarters []Card `json:"available_quarters"`
	}

	EventGameStarted struct {
		King *Player `json:"king"`
	}

	EventPickPhaseStarted struct {
		OpenLockedHeroes []Hero `json:"open_locked_heroes"`
		ClosedLockedHeroes int `json:"closed_locked_heroes"`
	}

	EventChooseHero struct {
		Heroes []Hero `json:"heroes"`
	}

	EventPlayerID struct {
		PlayerID PlayerID `json:"player_id"`
	}

	EventHeroSelected struct {
		Hero Hero `json:"hero"`
	}
)