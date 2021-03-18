package combat

import "go.uber.org/zap"

//eleType is a string representing an element i.e. HYDRO/PYRO/etc...
type eleType string

//ElementType should be pryo, Hydro, Cryo, Electro, Geo, Anemo and maybe dendro
const (
	Pyro     eleType = "pyro"
	Hydro    eleType = "hydro"
	Cryo     eleType = "cryo"
	Electro  eleType = "electro"
	Geo      eleType = "geo"
	Anemo    eleType = "anemo"
	Physical eleType = "physical"
	Frozen   eleType = "frozen"
)

//Enemy keeps track of the status of one enemy Enemy
type Enemy struct {
	Level  int64
	Resist float64

	//resist mods
	ResMod map[string]float64

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
func (e *Enemy) applyAura(ds snapshot) {
	//1A = 9.5s (570 frames) per unit, 2B = 6s (360 frames) per unit, 4C = 4.25s (255 frames) per unit
	//loop through existing auras and apply reactions if any
	if len(e.auras) > 1 {
		//this case should only happen with electro charge where there's 2 aura active at any one point in time
		for ele, a := range e.auras {
			if ele != ds.Element {
				zap.S().Debugw("apply aura", "aura", a, "existing ele", ele, "next ele", ds.Element)
			} else {
				zap.S().Debugf("not implemented!!!")
			}
		}
	} else if len(e.auras) == 1 {
		if a, ok := e.auras[ds.Element]; ok {
			next := aura{
				gauge:    ds.AuraGauge,
				unit:     a.unit,
				duration: auraDur(a.unit, ds.AuraGauge),
			}
			//refresh duration
			zap.S().Debugf("%v refreshed. unit: %v. new duration: %v", ds.Element, a.unit, next.duration)
			e.auras[ds.Element] = next
		} else {
			//apply reaction
			//The length of the freeze is based on the lowest remaining duration of the two elements applied.
			zap.S().Debugf("not implemented!!!")
		}
	} else {
		next := aura{
			gauge:    ds.AuraGauge,
			unit:     ds.AuraUnit,
			duration: auraDur(ds.AuraUnit, ds.AuraGauge),
		}
		zap.S().Debugf("%v applied (new). unit: %v. duration: %v", ds.Element, next.unit, next.duration)
		e.auras[ds.Element] = next
	}
}

func (e *Enemy) tick(s *Sim) {
	//tick down buffs and debuffs
	for k, v := range e.status {
		if v == 0 {
			delete(e.status, k)
		} else {
			e.status[k]--
		}
	}
	//tick down aura
	for k, v := range e.auras {
		if v.duration == 0 {
			print(s.Frame, true, "aura %v expired", k)
			delete(e.auras, k)
		} else {
			a := e.auras[k]
			a.duration--
			e.auras[k] = a
		}
	}
}
