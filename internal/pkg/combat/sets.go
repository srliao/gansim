package combat

import "go.uber.org/zap"

func setBlizzardStrayer(c *Character, s *Sim, count int) {
	if count >= 2 {
		c.Mods["Blizzard Strayer 2PC"] = make(map[StatType]float64)
		c.Mods["Blizzard Strayer 2PC"][CryoP] = 0.15
	}
	if count >= 4 {
		s.addEffect(func(snap *snapshot) bool {
			if snap.CharName != c.Profile.Name {
				return false
			}

			if _, ok := s.Target.auras[frozen]; ok {
				zap.S().Debugf("applying blizzard strayer 4pc buff on frozen target")
				snap.Stats[CR] += .4
			} else if _, ok := s.Target.auras[Cryo]; ok {
				zap.S().Debugf("applying blizzard strayer 4pc buff on cryo target")
				snap.Stats[CR] += .2
			}

			return false
		}, "blizzard strayer 4pc", preDamageHook)
	}
	//add flat stat to char
}
