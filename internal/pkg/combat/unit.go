package combat

import "go.uber.org/zap"

//eleType is a string representing an element i.e. HYDRO/PYRO/etc...
type eleType string

//ElementType should be pryo, hydro, cryo, electro, geo, anemo and maybe dendro
const (
	pyro     eleType = "pyro"
	hydro    eleType = "hydro"
	cryo     eleType = "cryo"
	electro  eleType = "electro"
	geo      eleType = "geo"
	anemo    eleType = "anemo"
	physical eleType = "physical"
	frozen   eleType = "frozen"
)

//unit keeps track of the status of one enemy unit
type unit struct {
	Level  int64
	Resist float64

	//tracking
	auras  map[eleType]aura
	status map[string]int //countdown to how long status last

	//stats
	damage float64 //total damage received
}

type aura struct {
	gauge    float64
	unit     string
	duration int
}

func auraDur(unit string, gauge float64) int {
	switch unit {
	case "A":
		return int(gauge * 9.5 * 60)
	case "B":
		return int(gauge * 6 * 60)
	case "C":
		return int(gauge * 4.25 * 60)
	}
	return 0
}

//applyAura applies an aura to the Unit, can trigger damage for superconduct, electrocharged, etc..
func (u *unit) applyAura(ds snapshot) {
	//1A = 9.5s (570 frames) per unit, 2B = 6s (360 frames) per unit, 4C = 4.25s (255 frames) per unit
	//loop through existing auras and apply reactions if any
	if len(u.auras) > 1 {
		//this case should only happen with electro charge where there's 2 aura active at any one point in time
		for e, a := range u.auras {
			if e != ds.element {
				zap.S().Debugw("apply aura", "aura", a, "existing ele", e, "next ele", ds.element)
			} else {
				zap.S().Debugf("not implemented!!!")
			}
		}
	} else if len(u.auras) == 1 {
		if a, ok := u.auras[ds.element]; ok {
			next := aura{
				gauge:    ds.auraGauge,
				unit:     a.unit,
				duration: auraDur(a.unit, ds.auraGauge),
			}
			//refresh duration
			zap.S().Debugf("%v refreshed. unit: %v. new duration: %v", ds.element, a.unit, next.duration)
			u.auras[ds.element] = next
		} else {
			//apply reaction
			//The length of the freeze is based on the lowest remaining duration of the two elements applied.
			zap.S().Debugf("not implemented!!!")
		}
	} else {
		next := aura{
			gauge:    ds.auraGauge,
			unit:     ds.auraUnit,
			duration: auraDur(ds.auraUnit, ds.auraGauge),
		}
		zap.S().Debugf("%v applied (new). unit: %v. duration: %v", ds.element, next.unit, next.duration)
		u.auras[ds.element] = next
	}
}

func (u *unit) tick(s *Sim) {
	//tick down buffs and debuffs
	for k, v := range u.status {
		if v == 0 {
			delete(u.status, k)
		} else {
			u.status[k]--
		}
	}
	//tick down aura
	for k, v := range u.auras {
		if v.duration == 0 {
			print(s.frame, true, "aura %v expired", k)
			delete(u.auras, k)
		} else {
			a := u.auras[k]
			a.duration--
			u.auras[k] = a
		}
	}
}
