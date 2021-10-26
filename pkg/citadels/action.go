package citadels

type CastFunc func(t *Table, caster *Player, ev Event) error

type Action struct {
	Event *Event

	do CastFunc
}

func NewAction(ev *Event, do CastFunc) Action {
	return Action{Event: ev, do: do}
}
