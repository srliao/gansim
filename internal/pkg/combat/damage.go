package combat

import (
	"math/rand"

	"go.uber.org/zap"
)

func (s *Sim) applyDamage(ds snapshot) float64 {

	ds.targetLvl = s.target.Level
	ds.targetRes = s.target.Resist

	for k, f := range s.effects[preDamageHook] {
		if f(&ds) {
			print(s.frame, true, "effect (pre damage) %v expired", k)
			delete(s.effects[preDamageHook], k)
		}
	}

	print(s.frame, true, "%v - %v triggered dmg", ds.char, ds.abil)

	damage := calcDmg(ds)

	for k, f := range s.effects[postDamageHook] {
		if f(&ds) {
			print(s.frame, true, "effect (post damage) %v expired", k)
			delete(s.effects[postDamageHook], k)
		}
	}

	s.target.damage += damage

	//apply aura
	if ds.applyAura {
		s.target.applyAura(ds)
	}

	return damage
}

type snapshot struct {
	char     string     //name of the character triggering the damage
	abil     string     //name of ability triggering the damage
	abilType ActionType //type of ability triggering the damage

	hitWeakPoint bool

	targetLvl int64
	targetRes float64

	mult         float64 //ability multiplier. could set to 0 from initial Mona dmg
	element      eleType //element of ability
	auraGauge    float64 //1 2 or 4
	auraUnit     string  //A, B, or C
	auraDuration int     //duration of the aura in units
	applyAura    bool    //if aura should be applied; false if under ICD
	useDef       bool    //default false
	flatDmg      float64 //flat dmg; so far only zhongli
	otherMult    float64 //so far just for xingqiu C4

	stats      map[StatType]float64 //total character stats including from artifact, bonuses, etc...
	baseAtk    float64              //base attack used in calc
	baseDef    float64              //base def used in calc
	dmgBonus   float64              //total damage bonus, including appropriate ele%, etc..
	charLvl    int64
	defMod     float64
	resMod     float64
	reactBonus float64 //reaction bonus %+ such as witch
}

func calcDmg(d snapshot) float64 {

	var st StatType
	switch d.element {
	case anemo:
		st = AnemoP
	case cryo:
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
	d.dmgBonus += d.stats[st]

	zap.S().Debugw("calc", "base atk", d.baseAtk, "flat +", d.stats[ATK], "% +", d.stats[ATKP], "bonus dmg", d.dmgBonus, "mul", d.mult)
	//calculate attack or def
	var a float64
	if d.useDef {
		a = d.baseDef*(1+d.stats[DEFP]) + d.stats[DEF]
	} else {
		a = d.baseAtk*(1+d.stats[ATKP]) + d.stats[ATK]
	}

	base := d.mult*a + d.flatDmg
	damage := base * (1 + d.dmgBonus)

	zap.S().Debugw("calc", "total atk", a, "base dmg", base, "dmg + bonus", damage)

	//make sure 0 <= cr <= 1
	if d.stats[CR] < 0 {
		d.stats[CR] = 0
	}
	if d.stats[CR] > 1 {
		d.stats[CR] = 1
	}

	//check if crit
	if rand.Float64() <= d.stats[CR] || d.hitWeakPoint {
		zap.S().Debugf("damage is crit!")
		damage = damage * (1 + d.stats[CD])
	}

	zap.S().Debugw("calc", "cr", d.stats[CR], "cd", d.stats[CD], "def adj", d.defMod, "res adj", d.resMod, "char lvl", d.charLvl, "target lvl", d.targetLvl)

	defmod := float64(d.charLvl+100) / (float64(d.charLvl+100) + float64(d.targetLvl+100)*(1-d.defMod))
	//apply def mod
	damage = damage * defmod
	//apply resist mod
	res := d.targetRes + d.resMod
	resmod := 1 - res/2
	if res >= 0 && res < 0.75 {
		resmod = 1 - res
	} else if res > 0.75 {
		resmod = 1 / (4*res + 1)
	}
	damage = damage * resmod
	zap.S().Debugw("calc", "def mod", defmod, "res mod", resmod)

	//apply other multiplier bonus
	if d.otherMult > 0 {
		damage = damage * d.otherMult
	}

	return damage
}
