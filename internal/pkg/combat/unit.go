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
	Buffs   map[string]bool
	Debuffs map[string]bool

	//action to be applied against this unit. action is applied when
	//delay = 0; otherwise decrement delay by 1
	//once applied, action should be removed from this list
	Actions []UnitAction

	//Hooks are called each tick; not sure usage yet
	Hooks map[string]func()

	//special hooks
	DamageHooks map[string]HookFunc //this is for actions that triggers on damage, such as fischl oz thinggy

	//stats
	Damage float64 //total damage received
}

//UnitAction represents an action to be performed on a unit after some delay (could be 0)
type UnitAction struct {
	//each action will either apply an aura, deal damage, apply a debuff, or apply a buff
	Callback HookFunc
	Delay    int
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

func (u *Unit) tick() {
	//check if any actions to apply
	var keep []UnitAction
	for i, a := range u.Actions {
		if a.Delay == 0 {
			//apply the action
			a.Callback(u)
		} else {
			u.Actions[i].Delay--
			keep = append(keep, u.Actions[i])
		}
	}
	u.Actions = keep
}

//DamageProfile describe the stats necessary to calculate the damage
type DamageProfile struct {
	Multiplier     float64              //ability multiplier. could set to 0 from initial Mona dmg
	Element        ElementType          //element of ability
	ApplyAura      bool                 //if aura should be applied; false if under ICD
	Stats          map[StatType]float64 //total character stats including from artifact, bonuses, etc...
	BaseAtk        float64              //base attack used in calc
	BaseDef        float64              //base def used in calc
	UseDef         bool                 //default false
	FlatDmg        float64              //flat dmg; so far only zhongli
	DmgBonus       float64              //total damage bonus, including appropriate ele%, etc..
	CharacterLevel int64
	DefMod         float64
	ResistMod      float64
	ReactionBonus  float64 //reaction bonus %+ such as witch
	OtherMult      float64 //so far just for xingqiu C4
}
