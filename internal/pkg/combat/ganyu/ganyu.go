package ganyu

import "github.com/srliao/gansim/internal/pkg/combat"

//New creates a character with Ganyu's abilities
func New() *combat.Character {
	var c combat.Character

	c.ChargeAttack = charge

	return &c
}

func charge(c *combat.Character, u *combat.Unit) int {

	travel := 100

	initial := func(u *combat.Unit) {
		//apply damage
		u.ApplyDamage()

		//if not ICD, apply aura
		if _, ok := c.ICD[combat.ActionTypeChargedAttack]; !ok {
			//not in map so no cd, so we can apply aura
			u.ApplyAura()
		}
	}
	//apply initial damage w/ travel time
	u.Actions = append(u.Actions, combat.UnitAction{
		Callback: initial,
		Delay:    travel, //fake travel time
	})

	//apply second bloom w/ more travel time
	bloom := func(u *combat.Unit) {
		//apply damage
		u.ApplyDamage()

		//if not ICD, apply aura
		if _, ok := c.ICD[combat.ActionTypeChargedAttack]; !ok {
			//not in map so no cd, so we can apply aura
			u.ApplyAura()
		}
	}

	u.Actions = append(u.Actions, combat.UnitAction{
		Callback: bloom,
		Delay:    travel + 30, //fake bloom travel time
	})

	//return animation cd
	return 0
}
