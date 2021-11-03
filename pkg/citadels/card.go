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
	Name  string
	Type  string
	Price int
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

func Witch() Hero {
	return Hero{
		Name:  "Witch",
		Turn:  1,
		Skill: Skill{
			Type: SkillTypeAnytime,
			Do: func(t *Table, caster *Player, ev Event) error {
				return nil
			},
		},
	}
}

func Blackmailer() Hero {
	return Hero{
		Name:  "Blackmailer",
		Turn:  2,
		Skill: Skill{
			Type: SkillTypeAnytime,
			Do: func(t *Table, caster *Player, ev Event) error {
				return nil
			},
		},
	}
}

func Enchantress() Hero {
	return Hero{
		Name:  "Enchantress",
		Turn:  3,
		Skill: Skill{
			Type: SkillTypeAnytime,
			Do: func(t *Table, caster *Player, ev Event) error {
				return nil
			},
		},
	}
}

func Abat() Hero {
	return Hero{
		Name:  "Abat",
		Turn:  5,
		Skill: Skill{
			Type: SkillTypeAnytime,
			Do: func(t *Table, caster *Player, ev Event) error {
				return nil
			},
		},
	}
}

func Alchemist() Hero {
	return Hero{
		Name:  "Alchemist",
		Turn:  6,
		Skill: Skill{
			Type: SkillTypeAnytime,
			Do: func(t *Table, caster *Player, ev Event) error {
				return nil
			},
		},
	}
}

func Architect() Hero {
	return Hero{
		Name:  "Architect",
		Turn:  7,
		Skill: Skill{
			Type: SkillTypeAnytime,
			Do: func(t *Table, caster *Player, ev Event) error {
				return nil
			},
		},
	}
}

func Warlord() Hero {
	return Hero{
		Name:  "Warlord",
		Turn:  8,
		Skill: Skill{
			Type: SkillTypeAnytime,
			Do: func(t *Table, caster *Player, ev Event) error {
				return nil
			},
		},
	}
}


func CustomsOfficer() Hero {
	return Hero{
		Name:  "CustomsOfficer",
		Turn:  9,
		Skill: Skill{
			Type: SkillTypeAnytime,
			Do: func(t *Table, caster *Player, ev Event) error {
				return nil
			},
		},
	}
}
