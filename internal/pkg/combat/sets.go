package combat

import "go.uber.org/zap"

func setBlizzardStrayer(c *character, s *Sim, count int) {
	if count >= 2 {
		c.statMods["Blizzard Strayer 2PC"] = make(map[StatType]float64)
		c.statMods["Blizzard Strayer 2PC"][CryoP] = 0.15
	}
	if count >= 4 {
		s.addEffect(func(snap *snapshot) bool {
			if snap.char != c.Name {
				return false
			}
			zap.S().Debugf("applying blizzard strayer 4pc buff")
			//check aura, if cryo crit + 20%
			if _, ok := s.target.auras[eTypeCryo]; ok {
				snap.stats[CR] += .2
			}
			//if frozen crit + 20%
			//TODO HOW TO CHECK IF FROZEN

			return false
		}, "blizzard strayer 4pc", preDamageHook)
	}
	//add flat stat to char
}
