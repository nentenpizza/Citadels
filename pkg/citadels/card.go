package citadels

import "encoding/json"

// Quarter Types
const (
	QuarterTypeNoble = "Noble"
	QuarterTypeMilitary = "Military"
	QuarterTypeTrade = "Trade"
	QuarterTypeSpiritual = "Spiritual"
	QuarterTypeSpecial = "Special"
)

type Quarter struct {
	Name string
	Type string
	Cost int
}

// Turns for heroes
const (
	EmperorTurn = 4
	EnchantressTurn = 3
)

// Skill types
const (
	SkillTypeAnytime = "skill.type.anytime"
	SkillTypeAtStart = "skill.type.at.start"
)

type Skill struct {
	Type string
	Do CastFunc
}

type Hero struct {
	// Name of hero
	Name string `json:"name"`

	// Turn is hero's move number
	// heroes makes moves in specific order from 1 to 9
	Turn int `json:"turn"`

	Skill Skill
}

func Emperor() Hero {
	return Hero{
		Name:  "Emperor",
		Turn:  4,
		Skill: Skill{
			Type: SkillTypeAnytime,
			Do: func(t *Table, caster *Player, ev Event) error {
				b, err := json.Marshal(ev.Data)
				if err != nil {
					return err
				}
				var e EventEmperorSkill
				err = json.Unmarshal(b, &e)
				if err != nil {
					return ErrWrongEventData
				}

				target, ok := t.PlayerByID(string(e.TargetID))
				if !ok {
					return ErrPlayerNotExists
				}

				if target.ID == caster.ID {
					return ErrCannotCastOnMyself
				}

				if e.Coin {
					if target.Coins > 0 {
						target.giveCoins(caster, 1)
						return nil
					}
					ev.Error = ErrorTypeTargetHasNoCoins
					caster.Notify(ev)
				}
				if !e.Coin {
					if len(target.AvailableQuarters) > 0 {
						target.giveRandomCards(caster, 1)

						return nil
					}
					ev.Error = ErrorTypeTargetHasNoCards
					caster.Notify(ev)
				}

				return nil
			},
		},
	}
}