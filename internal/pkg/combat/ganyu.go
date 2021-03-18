package combat

import (
	"fmt"
	"log"

	"go.uber.org/zap"
)

func newGanyu(c *character) {
	c.chargeAttack = ganyuCAFunc
	c.burst = ganyuBurstFunc
	c.skill = ganyuSkillFunc
	c.attack = ganyuAttackFunc
	c.plungeAttack = ganyuPlungeFunc
	//start with full energy
	c.maxEnergy = 60
	c.maxStamina = 240
	c.energy = 60
	c.stamina = 240
}

func ganyuCAFunc(s *Sim) int {
	c := s.characters[s.active]
	//if it's the wrong char somehow (shouldn't be), exit
	if c.profile.Name != "Ganyu" {
		log.Panic("wrong character executing ability")
	}

	i := 0
	initial := func(s *Sim) bool {
		if i < 20 {
			i++
			return false
		}
		//abil
		d := c.snapshot(cryo)
		d.abil = "Frost Flake Arrow"
		d.abilType = ActionTypeChargedAttack
		d.hitWeakPoint = true
		d.mult = 2.304
		d.auraGauge = 1
		d.auraUnit = "A"
		d.applyAura = true
		//if not ICD, apply aura
		if _, ok := c.cooldown[storekey("ICD", "charge")]; !ok {
			d.applyAura = true
		}
		//check if A4 talent is
		if _, ok := c.cooldown["A2"]; ok {
			d.stats[CR] += 0.2
		}
		c.cooldown["A2"] = 5 * 60
		//apply damage
		damage := s.applyDamage(d)
		zap.S().Infof("[%v]: Ganyu frost arrow dealt %.0f damage", fts(s.frame), damage)
		return true
	}

	b := 0
	//apply second bloom w/ more travel time
	bloom := func(s *Sim) bool {
		if b < 50 {
			b++
			return false
		}
		//abil
		d := c.snapshot(cryo)
		d.abil = "Frost Flake Bloom"
		d.abilType = ActionTypeChargedAttack
		d.mult = 3.9168
		d.applyAura = true
		d.auraGauge = 1
		d.auraUnit = "A"
		//if not ICD, apply aura
		if _, ok := c.cooldown[storekey("ICD", "charge")]; !ok {
			d.applyAura = true
		}
		if _, ok := c.cooldown["A2"]; ok {
			d.stats[CR] += 0.2
		}
		//apply damage
		damage := s.applyDamage(d)
		zap.S().Infof("[%v]: Ganyu frost flake bloom dealt %.0f damage", fts(s.frame), damage)
		return true
	}
	s.addAction(initial, fmt.Sprintf("%v-Ganyu-CA-FFA", s.frame))
	s.addAction(bloom, fmt.Sprintf("%v-Ganyu-CA-FFB", s.frame))

	//return animation cd
	return 137
}

func ganyuAttackFunc(s *Sim) int {
	return 0
}

func ganyuPlungeFunc(s *Sim) int {
	return 0
}

func ganyuBurstFunc(s *Sim) int {
	current := s.characters[s.active]
	//if it's the wrong char somehow (shouldn't be), exit
	if current.profile.Name != "Ganyu" {
		log.Panic("wrong character executing ability")
	}

	//snap shot stats at cast time here
	d := current.snapshot(cryo)
	d.abil = "Celestial Shower"
	d.abilType = ActionTypeBurst
	d.mult = 0.938
	d.applyAura = true
	d.auraGauge = 1
	d.auraUnit = "A"
	d.auraDuration = 570 //9.5s * 60 frames

	//apply weapon stats here
	//burst should be instant
	//should add a hook to the unit, triggering damage every 1 sec
	//also add a field effect
	tick := 0
	storm := func(s *Sim) bool {
		if tick > 900 {
			return true
		}
		//check if multiples of 60s; also add an initial delay of 120 frames
		if tick%60 != 0 || tick < 120 {
			tick++
			return false
		}
		//do damage
		damage := s.applyDamage(d)
		zap.S().Infof("[%v]: Ganyu burst (tick) dealt %.0f damage", fts(s.frame), damage)
		tick++
		return false
	}
	s.addAction(storm, fmt.Sprintf("%v-Ganyu-Burst", s.frame))
	//add cooldown to sim
	current.cooldown[storekey("cd", "burst")] = 15 * 60

	return 122
}

func ganyuSkillFunc(s *Sim) int {
	current := s.characters[s.active]
	//if it's the wrong char somehow (shouldn't be), exit
	if current.profile.Name != "Ganyu" {
		log.Panic("wrong character executing ability")
	}
	//snap shot stats at cast time here
	d := current.snapshot(cryo)
	d.mult = 1.848
	d.applyAura = true
	d.auraGauge = 1
	d.auraUnit = "A"
	d.auraDuration = 570 //9.5s * 60 frames

	tick := 0
	flower := func(s *Sim) bool {
		if tick < 6*60 {
			return false
		}
		//do damage
		damage := s.applyDamage(d)
		zap.S().Infof("[%v]: Ganyu ice lotus (tick) dealt %.0f damage", fts(s.frame), damage)
		tick++
		return false
	}
	s.addAction(flower, fmt.Sprintf("%v-Ganyu-Skill", s.frame))
	//add cooldown to sim
	current.cooldown[storekey("cd", "burst")] = 15 * 60

	return 30
}
