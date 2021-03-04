package combat

import "fmt"

//NewGanyu creates a character with Ganyu's abilities
func NewGanyu() *Character {
	var c Character
	c.Name = "Ganyu"

	c.ChargeAttack = func(s *Sim) int {

		travel := 50
		c := s.Actors[s.Active]
		u := s.Target

		//snap shot stats at cast time here
		var d DamageProfile
		for k, v := range c.ArtifactStats {
			d.Stats[k] = v
		}
		//other stats
		d.BaseAtk = c.BaseAtk
		d.CharacterLevel = c.Level
		d.Element = ElementTypeCryo

		//ganyu ascension bonuses

		//ganyu c1
		d.DefMod = -0.1

		//add dmg bonus for cryo
		d.DmgBonus += c.ArtifactStats[HydroP]

		initial := func(u *Unit) {
			//abil
			d.Multiplier = 1
			//if not ICD, apply aura
			if _, ok := c.ICD[ActionTypeChargedAttack]; !ok {
				d.ApplyAura = true
			}
			//apply damage
			damage := u.ApplyDamage(d)
			fmt.Printf("[%v] Ganyu frost arrow dealt %.0f damage\n", s.Frame, damage)
		}
		//apply initial damage w/ travel time
		u.Actions = append(u.Actions, UnitAction{
			Callback: initial,
			Delay:    travel, //fake travel time
		})

		//apply second bloom w/ more travel time
		bloom := func(u *Unit) {
			//abil
			d.Multiplier = 1
			//if not ICD, apply aura
			if _, ok := c.ICD[ActionTypeChargedAttack]; !ok {
				d.ApplyAura = true
			}
			//apply damage
			damage := u.ApplyDamage(d)
			fmt.Printf("[%v] Ganyu frost flake bloom dealt %.0f damage\n", s.Frame, damage)
		}

		u.Actions = append(u.Actions, UnitAction{
			Callback: bloom,
			Delay:    travel + 30, //fake bloom travel time
		})

		//return animation cd
		return 90
	}

	return &c
}
