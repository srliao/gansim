package combat

import (
	"math/rand"
)

//Unit keeps track of the status of one enemy Unit
type Unit struct {
	Level  int64
	Resist float64

	//Auras
	Auras   map[ElementType]Aura
	Buffs   map[string]int //countdown to how long status last
	Debuffs map[string]int //countdown to how long statu last

	//stats
	Damage float64 //total damage received
}

//ApplyDamage applies ability damage to an Unit
func (u *Unit) ApplyDamage(d DamageProfile) float64 {
	//calculate attack or def
	var a float64
	if d.UseDef {
		a = d.BaseDef*(1+d.Stats[DEFP]) + d.Stats[DEF]
	} else {
		a = d.BaseAtk*(1+d.Stats[ATKP]) + d.Stats[ATK]
	}
	base := d.Multiplier*a + d.FlatDmg

	damage := base * (1 + d.DmgBonus)

	//check if crit
	if rand.Float64() <= d.Stats[CR] {
		damage = damage * (1 + d.Stats[CD])
	}

	defmod := float64(d.CharacterLevel+100) / (float64(d.CharacterLevel+100) + float64(u.Level+100)*(1-d.DefMod))
	//apply def mod
	damage = damage * defmod
	//apply resist mod
	res := u.Resist + d.ResistMod
	resmod := 1 - res/2
	if res >= 0 && res < 0.75 {
		resmod = 1 - res
	} else if res > 0.75 {
		resmod = 1 / (4*res + 1)
	}
	damage = damage * resmod

	//apply amp mod - TODO
	if d.ApplyAura {

	}

	//apply other multiplier bonus
	if d.OtherMult > 0 {
		damage = damage * d.OtherMult
	}

	u.Damage += damage

	return damage
}

//ApplyAura applies an aura to the Unit
func (u *Unit) ApplyAura() {
	//can trigger apply damage for superconduct, electrocharged, etc..
}

func (u *Unit) tick(s *Sim) {
	//tick down buffs and debuffs
	for k, v := range u.Buffs {
		if v == 0 {
			delete(u.Buffs, k)
		} else {
			u.Buffs[k]--
		}
	}
	for k, v := range u.Debuffs {
		if v == 0 {
			delete(u.Debuffs, k)
		} else {
			u.Debuffs[k]--
		}
	}
	//tick down aura
}

//DamageProfile describe the stats necessary to calculate the damage
type DamageProfile struct {
	Multiplier float64     //ability multiplier. could set to 0 from initial Mona dmg
	Element    ElementType //element of ability
	ApplyAura  bool        //if aura should be applied; false if under ICD
	UseDef     bool        //default false
	FlatDmg    float64     //flat dmg; so far only zhongli
	OtherMult  float64     //so far just for xingqiu C4

	Stats          map[StatType]float64 //total character stats including from artifact, bonuses, etc...
	BaseAtk        float64              //base attack used in calc
	BaseDef        float64              //base def used in calc
	DmgBonus       float64              //total damage bonus, including appropriate ele%, etc..
	CharacterLevel int64
	DefMod         float64
	ResistMod      float64
	ReactionBonus  float64 //reaction bonus %+ such as witch
}
