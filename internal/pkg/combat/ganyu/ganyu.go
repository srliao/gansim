package ganyu

import (
	"fmt"

	"github.com/srliao/gansim/internal/pkg/combat"
)

//New creates a character with Ganyu's abilities
func New() *combat.Character {
	var c combat.Character

	c.ChargeAttack = charge

	return &c
}

func charge(s *combat.Sim) int {

	travel := 100
	c := s.Actors[s.Active]
	u := s.Target

	//snap shot stats at cast time here
	var d combat.DamageProfile
	for k, v := range c.ArtifactStats {
		d.Stats[k] = v
	}
	//other stats
	d.BaseAtk = c.BaseAtk
	d.CharacterLevel = c.Level
	d.Element = combat.ElementTypeCryo

	//ganyu ascension bonuses

	//ganyu c1
	d.DefMod = -0.1

	//add dmg bonus for cryo
	d.DmgBonus += c.ArtifactStats[combat.HydroP]

	initial := func(u *combat.Unit) {
		//abil
		d.Multiplier = 1
		//if not ICD, apply aura
		if _, ok := c.ICD[combat.ActionTypeChargedAttack]; !ok {
			d.ApplyAura = true
		}
		//apply damage
		damage := u.ApplyDamage(d)
		fmt.Printf("[%v] Ganyu frost arrow dealt %.0f damage\n", s.Frame, damage)
	}
	//apply initial damage w/ travel time
	u.Actions = append(u.Actions, combat.UnitAction{
		Callback: initial,
		Delay:    travel, //fake travel time
	})

	//apply second bloom w/ more travel time
	bloom := func(u *combat.Unit) {
		//abil
		d.Multiplier = 1
		//if not ICD, apply aura
		if _, ok := c.ICD[combat.ActionTypeChargedAttack]; !ok {
			d.ApplyAura = true
		}
		//apply damage
		damage := u.ApplyDamage(d)
		fmt.Printf("[%v] Ganyu frost flake bloom dealt %.0f damage\n", s.Frame, damage)
	}

	u.Actions = append(u.Actions, combat.UnitAction{
		Callback: bloom,
		Delay:    travel + 30, //fake bloom travel time
	})

	//return animation cd
	return 0
}
