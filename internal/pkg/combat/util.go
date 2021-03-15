package combat

import "fmt"

func print(f int, msg string, data ...interface{}) {
	fmt.Printf("[%.2fs|%v]: %v\n", float64(f)/60, f, fmt.Sprintf(msg, data...))
}

func dProfileBuilder(c *character) DamageProfile {
	var d DamageProfile
	d.Stats = make(map[StatType]float64)
	for k, v := range c.stats {
		d.Stats[k] = v
	}
	//other stats
	d.BaseAtk = c.BaseAtk + c.WeaponAtk
	d.CharacterLevel = c.Level
	d.BaseDef = c.BaseDef

	return d
}

func storekey(t, k string) string {
	return fmt.Sprintf("%v-%v", t, k)
}

/**
Stats          map[StatType]float64 //total character stats including from artifact, bonuses, etc...
BaseAtk        float64              //base attack used in calc
BaseDef        float64              //base def used in calc
CharacterLevel int64

DmgBonus       float64              //total damage bonus, including appropriate ele%, etc..
DefMod         float64
ResistMod      float64
ReactionBonus  float64 //reaction bonus %+ such as witch
**/
