package citadels

import "math/rand"

// BroadcastEvent sends Event to all players at the table
func (t *Table) BroadcastEvent(e Event) {
	t.Lock()
	defer t.Unlock()
	t.doBroadcastEvent(e)
}

func (t *Table) doBroadcastEvent(e Event) {
	for _, p := range t.players {
		p.updates <- e
	}
}

// PlayerByID returns player with id, if not exists returns nil, false
func (t *Table) PlayerByID(id string) (*Player, bool) {
	t.Lock()
	defer t.Unlock()
	return t.playerByID(id)
}

// PlayerByID returns player with id, if not exists returns nil, false
func (t *Table) playerByID(id string) (*Player, bool) {
	p, ok := t.players[PlayerID(id)]
	return p, ok
}

func (t *Table) forceSelecting() {
	randomIndex := rand.Intn(len(t.heroesToSelect))
	hero := t.heroesToSelect[randomIndex]
	t.selecting.Hero = hero
	t.selecting.Notify(Event{
		Type:  EventTypeHeroSelected,
		Data:  EventHeroSelected{Hero: hero},
	})
	t.heroesToSelect = removeHero(t.heroesToSelect, randomIndex)
	t.nextSelecting()
}

func removeQuarter(slice []Quarter, s int) []Quarter {
	return append(slice[:s], slice[s+1:]...)
}

func removeQuarterByName(slice []Quarter, name string) []Quarter {
	for i, q := range slice{
		if q.Name == name{
			return removeQuarter(slice, i)
		}
	}
	return slice
}

func removeHero(slice []Hero, s int) []Hero {
	return append(slice[:s], slice[s+1:]...)
}
