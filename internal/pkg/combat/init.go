package combat

var setBonus = make(map[string]setBonusFunc)

type setBonusFunc func(c *Character, s *Sim, count int)

func init() {
	setBonus["Blizzard Strayer"] = setBlizzardStrayer
}
