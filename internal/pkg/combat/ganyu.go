package combat

import (
	"fmt"
	"log"
)

func newGanyu() *character {
	var c character
	c.Name = "Ganyu"
	c.cooldown = make(map[string]int)
	c.store = make(map[string]interface{})
	c.tickHooks = make(map[string]HookFunc)
	c.stats = make(map[StatType]float64)
	c.chargeAttack = ganyuCAFunc
	c.burst = ganyuBurstFunc
	c.skill = ganyuSkillFunc
	c.attack = ganyuAttackFunc

	return &c
}

func ganyuCAFunc(s *Sim) int {
	current := s.characters[s.active]
	//if it's the wrong char somehow (shouldn't be), exit
	if current.Name != "Ganyu" {
		log.Panic("wrong character executing ability")
	}
	//snap shot stats at cast time here
	d := dProfileBuilder(current)
	d.Element = ElementTypeCryo
	d.DmgBonus = current.stats[CryoP]

	//apply weapon stats here

	i := 0
	initial := func(s *Sim) bool {
		if i < 20 {
			i++
			return false
		}
		//abil
		d.Multiplier = 2.304
		//if not ICD, apply aura
		if _, ok := current.cooldown[storekey("ICD", "charge")]; !ok {
			d.ApplyAura = true
		}
		//apply damage
		damage := s.handleDamage(d)
		print(s.frame, "Ganyu frost arrow dealt %.0f damage", damage)
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
		d.Multiplier = 3.9168
		//if not ICD, apply aura
		if _, ok := current.cooldown[storekey("ICD", "charge")]; !ok {
			d.ApplyAura = true
		}
		//apply damage
		damage := s.handleDamage(d)
		print(s.frame, "Ganyu frost flake bloom dealt %.0f damage", damage)
		return true
	}
	s.useEffect(initial, fmt.Sprintf("%v-Ganyu-CA-FFA", s.frame), actionHook)
	s.useEffect(bloom, fmt.Sprintf("%v-Ganyu-CA-FFB", s.frame), actionHook)

	//return animation cd
	return 137
}

func ganyuAttackFunc(s *Sim) int {

	return 0
}

func ganyuBurstFunc(s *Sim) int {
	current := s.characters[s.active]
	//if it's the wrong char somehow (shouldn't be), exit
	if current.Name != "Ganyu" {
		log.Panic("wrong character executing ability")
	}
	//snap shot stats at cast time here
	var d DamageProfile
	d.Stats = make(map[StatType]float64)
	for k, v := range current.stats {
		d.Stats[k] = v
	}
	d.BaseAtk = current.BaseAtk + current.WeaponAtk
	d.CharacterLevel = current.Level
	d.BaseDef = current.BaseDef
	d.UseDef = false
	d.Element = ElementTypeHydro
	d.DmgBonus = current.stats[CryoP]

	//apply weapon stats here
	//burst should be instant
	//should add a hook to the unit, triggering damage every 1 sec
	//also add a field effect
	tick := 0
	storm := func(s *Sim) bool {
		d.Multiplier = 0.938
		if tick > 900 {
			return true
		}
		//check if multiples of 60s; also add an initial delay of 120 frames
		if tick%60 != 0 || tick < 120 {
			tick++
			return false
		}
		//do damage
		damage := s.handleDamage(d)
		print(s.frame, "\t [tick] Ganyu burst dealt %.0f damage", damage)
		tick++
		return false
	}
	s.useEffect(storm, fmt.Sprintf("%v-Ganyu-Burst", s.frame), actionHook)
	//add cooldown to sim
	current.cooldown[storekey("cd", "burst")] = 15 * 60

	return 122
}

func ganyuSkillFunc(s *Sim) int {
	current := s.characters[s.active]
	//if it's the wrong char somehow (shouldn't be), exit
	if current.Name != "Ganyu" {
		log.Panic("wrong character executing ability")
	}
	//snap shot stats at cast time here
	var d DamageProfile
	d.Stats = make(map[StatType]float64)
	for k, v := range current.stats {
		d.Stats[k] = v
	}
	d.BaseAtk = current.BaseAtk + current.WeaponAtk
	d.CharacterLevel = current.Level
	d.BaseDef = current.BaseDef
	d.UseDef = false
	d.Element = ElementTypeHydro
	d.DmgBonus = current.stats[CryoP]
	tick := 0
	flower := func(s *Sim) bool {
		d.Multiplier = 1.848
		if tick < 6*60 {
			return false
		}
		//do damage
		damage := s.handleDamage(d)
		print(s.frame, "\t [tick] Ganyu ice lotus dealt %.0f damage", damage)
		tick++
		return false
	}
	s.useEffect(flower, fmt.Sprintf("%v-Ganyu-Skill", s.frame), actionHook)
	//add cooldown to sim
	current.cooldown[storekey("cd", "burst")] = 15 * 60

	return 30
}
