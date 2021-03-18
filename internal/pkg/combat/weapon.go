package combat

import (
	"fmt"

	"go.uber.org/zap"
)

func weaponPrototypeCrescent(c *character, s *Sim, r int) {
	//add on hit effect to sim?
	s.addEffect(func(snap *snapshot) bool {
		//check if char is correct?
		if snap.char != c.Name {
			return false
		}
		//check if weakpoint triggered
		if !snap.hitWeakPoint {
			return false
		}
		//add a new action that adds % dmg to current char and removes itself after
		//10 seconds
		tick := 0
		s.addAction(func(s *Sim) bool {
			if tick >= 10*60 {
				delete(c.statMods, "Prototype-Crescent-Proc")
				zap.S().Debugw("prototype crescent buff expired", "tick", tick)
				return true
			}
			tick++
			if _, ok := c.statMods["Prototype-Crescent-Proc"]; !ok {
				c.statMods["Prototype-Crescent-Proc"] = make(map[StatType]float64)
				atkmod := 0.36
				switch r {
				case 2:
					atkmod = 0.45
				case 3:
					atkmod = 0.54
				case 4:
					atkmod = 0.63
				case 5:
					atkmod = 0.72
				}
				zap.S().Debugw("applying prototype crescent buff", "%", atkmod, "tick", tick)
				c.statMods["Prototype-Crescent-Proc"][ATKP] = atkmod
			}
			return false
		}, fmt.Sprintf("%v-Prototype-Crescent-Proc-%v", s.frame, c.Name))
		return false
	}, "prototype-crescent-proc", postDamageHook)
}
