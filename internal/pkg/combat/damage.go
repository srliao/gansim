package combat

import (
	"math/rand"

	"go.uber.org/zap"
)

func (s *Sim) ApplyDamage(ds snapshot) float64 {

	ds.TargetLvl = s.Target.Level
	ds.TargetRes = s.Target.Resist

	for k, f := range s.effects[preDamageHook] {
		if f(&ds) {
			print(s.Frame, true, "effect (pre damage) %v expired", k)
			delete(s.effects[preDamageHook], k)
		}
	}

	print(s.Frame, true, "%v - %v triggered dmg", ds.CharName, ds.Abil)

	damage := calcDmg(ds)

	for k, f := range s.effects[postDamageHook] {
		if f(&ds) {
			print(s.Frame, true, "effect (post damage) %v expired", k)
			delete(s.effects[postDamageHook], k)
		}
	}

	s.Target.damage += damage

	//apply aura
	if ds.ApplyAura {
		s.Target.applyAura(ds)
	}

	return damage
}

type snapshot struct {
	CharName string     //name of the character triggering the damage
	Abil     string     //name of ability triggering the damage
	AbilType ActionType //type of ability triggering the damage

	HitWeakPoint bool

	TargetLvl int64
	TargetRes float64

	Mult         float64 //ability multiplier. could set to 0 from initial Mona dmg
	Element      eleType //element of ability
	AuraGauge    float64 //1 2 or 4
	AuraUnit     string  //A, B, or C
	AuraDuration int     //duration of the aura in units
	ApplyAura    bool    //if aura should be applied; false if under ICD
	UseDef       bool    //default false
	FlatDmg      float64 //flat dmg; so far only zhongli
	OtherMult    float64 //so far just for xingqiu C4

	Stats      map[StatType]float64 //total character stats including from artifact, bonuses, etc...
	BaseAtk    float64              //base attack used in calc
	BaseDef    float64              //base def used in calc
	DmgBonus   float64              //total damage bonus, including appropriate ele%, etc..
	CharLvl    int64
	DefMod     float64
	ResMod     float64
	ReactBonus float64 //reaction bonus %+ such as witch
}

func calcDmg(d snapshot) float64 {

	var st StatType
	switch d.Element {
	case anemo:
		st = AnemoP
	case Cryo:
		st = CryoP
	case electro:
		st = ElectroP
	case geo:
		st = GeoP
	case hydro:
		st = HydroP
	case pyro:
		st = PyroP
	case physical:
		st = PhyP
	}
	d.DmgBonus += d.Stats[st]

	zap.S().Debugw("calc", "base atk", d.BaseAtk, "flat +", d.Stats[ATK], "% +", d.Stats[ATKP], "bonus dmg", d.DmgBonus, "mul", d.Mult)
	//calculate attack or def
	var a float64
	if d.UseDef {
		a = d.BaseDef*(1+d.Stats[DEFP]) + d.Stats[DEF]
	} else {
		a = d.BaseAtk*(1+d.Stats[ATKP]) + d.Stats[ATK]
	}

	base := d.Mult*a + d.FlatDmg
	damage := base * (1 + d.DmgBonus)

	zap.S().Debugw("calc", "total atk", a, "base dmg", base, "dmg + bonus", damage)

	//make sure 0 <= cr <= 1
	if d.Stats[CR] < 0 {
		d.Stats[CR] = 0
	}
	if d.Stats[CR] > 1 {
		d.Stats[CR] = 1
	}

	//check if crit
	if rand.Float64() <= d.Stats[CR] || d.HitWeakPoint {
		zap.S().Debugf("damage is crit!")
		damage = damage * (1 + d.Stats[CD])
	}

	zap.S().Debugw("calc", "cr", d.Stats[CR], "cd", d.Stats[CD], "def adj", d.DefMod, "res adj", d.ResMod, "char lvl", d.CharLvl, "target lvl", d.TargetLvl)

	defmod := float64(d.CharLvl+100) / (float64(d.CharLvl+100) + float64(d.TargetLvl+100)*(1-d.DefMod))
	//apply def mod
	damage = damage * defmod
	//apply resist mod
	res := d.TargetRes + d.ResMod
	resmod := 1 - res/2
	if res >= 0 && res < 0.75 {
		resmod = 1 - res
	} else if res > 0.75 {
		resmod = 1 / (4*res + 1)
	}
	damage = damage * resmod
	zap.S().Debugw("calc", "def mod", defmod, "res mod", resmod)

	//apply other multiplier bonus
	if d.OtherMult > 0 {
		damage = damage * d.OtherMult
	}

	return damage
}
