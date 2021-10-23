package citadels

// Event Types
var (
	EventTypePickPhaseStarted = "phase.pick.started"

	// Targeted events

	EventTypeSelectHero = "select.hero"
)

type Event struct {
	Type string
	Data interface{}
}

// Events Data
type(
	EventTargeted struct {
		TargetID string
	}

	e struct {

	}
)