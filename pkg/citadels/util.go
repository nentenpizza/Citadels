package citadels

// BroadcastEvent sends Event to all players at the table
func (t *Table) BroadcastEvent(e *Event) {
	for _, p := range t.players {
		p.updates <- e
	}
}