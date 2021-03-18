package combat

//unit keeps track of the status of one enemy unit
type unit struct {
	Level  int64
	Resist float64

	//tracking
	auras  map[eleType]int
	status map[string]int //countdown to how long status last

	//stats
	damage float64 //total damage received
}

//applyAura applies an aura to the Unit
func (u *unit) applyAura(ele eleType, dur int) {
	//can trigger apply damage for superconduct, electrocharged, etc..
	u.auras[ele] += dur
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
		if v == 0 {
			print(s.frame, true, "aura %v expired", k)
			delete(u.auras, k)
		} else {
			u.auras[k]--
		}
	}
}
