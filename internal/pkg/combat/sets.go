package combat

import "go.uber.org/zap"

func setBlizzardStrayer(c *character, s *Sim, count int) {
	if count >= 2 {
		c.statMods["Blizzard Strayer 2PC"] = make(map[StatType]float64)
		c.statMods["Blizzard Strayer 2PC"][CryoP] = 0.15
	}
	if count >= 4 {
		s.addEffect(func(snap *snapshot) bool {
			if snap.char != c.profile.Name {
				return false
			}

			if _, ok := s.target.auras[frozen]; ok {
				zap.S().Debugf("applying blizzard strayer 4pc buff on frozen target")
				snap.stats[CR] += .4
			} else if _, ok := s.target.auras[cryo]; ok {
				zap.S().Debugf("applying blizzard strayer 4pc buff on cryo target")
				snap.stats[CR] += .2
			}

			return false
		}, "blizzard strayer 4pc", preDamageHook)
	}
	//add flat stat to char
}
