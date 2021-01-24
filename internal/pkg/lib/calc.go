package lib

import (
	"fmt"
	"math"
)

//Profile describe a damage profile to calculate
type Profile struct {
	CharLevel     float64 `yaml:"CharacterLevel"`
	CharBaseAtk   float64 `yaml:"CharacterBaseAtk"`
	WeaponBaseAtk float64 `yaml:"WeaponBaseAtk"`
	EnemyLevel    float64 `yaml:"EnemyLevel"`
	//abilities
	Abilities []struct {
		Talent            float64   `default:"1.0" yaml:"Talent"`
		VapMeltMultiplier float64   `default:"1.0" yaml:"VaporizeOrMeltMultiplier"`
		AtkMod            []float64 `yaml:"AtkMod"`
		EleMod            []float64 `yaml:"EleMod"`
		CCMod             []float64 `yaml:"CCMod"`
		CDMod             []float64 `yaml:"CDMod"`
		DmgMod            []float64 `yaml:"DmgMod"`
		EMMod             []float64 `yaml:"EMMod"`
		ReactionBonus     []float64 `yaml:"ReactionBonus"`
		ResistMod         []float64 `yaml:"ResistMod"`
		DefShredMod       []float64 `yaml:"DefShredMod"`
	} `yaml:"Abilities"`
}

//DmgResult is result for one abil
type DmgResult struct {
	Normal float64
	Avg    float64
	Crit   float64
}

//Calc calculates normal, avg, and crit damage given a profile
func Calc(p Profile, s map[Slot]Artifact, showDebug bool) []DmgResult {

	artifactStats := make(map[StatType]float64)

	for _, a := range s {
		artifactStats[a.MainStat.Type] += a.MainStat.Value

		for _, v := range a.Substat {
			artifactStats[v.Type] += v.Value
		}

	}
	baseAtk := p.CharBaseAtk + p.WeaponBaseAtk

	if showDebug {
		fmt.Print("--------------------------\n")
		fmt.Printf("\tCharacter + weapon base atk: %.2f\n", baseAtk)
		fmt.Printf("\t%v\n", artifactStats)
		// fmt.Printf("\t%v\n", artifactStats[ATK])
		// fmt.Printf("\t%v\n", artifactStats[ATKP])
		// fmt.Printf("\t%v\n", artifactStats[CR])
		// fmt.Printf("\t%v\n", artifactStats[CD])
		// fmt.Printf("\t%v\n", artifactStats[EleP])
		// fmt.Printf("\t%v\n", artifactStats[EM])
	}

	var r []DmgResult

	for _, ab := range p.Abilities {
		var totalAtk, atk, atkp, cr, cd, elep, em, dmgBonus, reactionBonus, defShred, defAdj, resAdj, res float64
		totalAtk = baseAtk

		atk += artifactStats[ATK]
		atkp += artifactStats[ATKP]
		cr += artifactStats[CR]
		cd += artifactStats[CD]
		elep += artifactStats[EleP]
		em += artifactStats[EM]

		if showDebug {
			fmt.Print("--------------------------\n")
			fmt.Print("\tArtifact stats -")
			fmt.Printf(" atk: %.2f", atk)
			fmt.Printf(" atkp: %.2f", atkp)
			fmt.Printf(" cc: %.2f", cr)
			fmt.Printf(" cd: %.2f", cd)
			fmt.Printf(" elep: %.2f", elep)
			fmt.Printf(" em: %.2f", em)
			fmt.Print("\n")
		}

		//add up atk % mods
		for _, v := range ab.AtkMod {
			atkp += v
		}

		totalAtk = totalAtk*(1+atkp) + atk

		//add up dmg mods
		for _, v := range ab.EleMod {
			dmgBonus += v
		}
		for _, v := range ab.DmgMod {
			dmgBonus += v
		}

		dmgBonus += elep //add in ele bonus from artifacts

		//add up crit mods
		for _, v := range ab.CCMod {
			cr += v
		}
		for _, v := range ab.CDMod {
			cd += v
		}

		//cap cc at 1
		if cr > 1 {
			cr = 1
		}
		if cr < 0 {
			cr = 0
		}
		if cd < 0 {
			cd = 0
		}

		//add up em mod
		for _, v := range ab.EMMod {
			em += v
		}

		//add up def shreds
		for _, v := range ab.DefShredMod {
			defShred += v
		}
		//calculate def adjustment
		defAdj = (100 + p.CharLevel) / ((100 + p.CharLevel) + (100+p.EnemyLevel)*(1-defShred))

		for _, v := range ab.ResistMod {
			res += v
		}

		if res < 0 {
			resAdj = 1 - (res / 2)
		} else if res < 0.75 {
			resAdj = 1 - res
		} else {
			resAdj = 1 / (4*res + 1)
		}

		for _, v := range ab.ReactionBonus {
			reactionBonus += v
		}

		vmMult := 1.0
		if ab.VapMeltMultiplier == 1.5 || ab.VapMeltMultiplier == 2 {
			vmMult = ab.VapMeltMultiplier*(1+0.00189266831*em*math.Exp(-0.000505*em)) + reactionBonus
		}

		if showDebug {
			fmt.Print("--------------------------\n")
			fmt.Print("\tAdjusted stats -")
			fmt.Printf(" total atk: %.4f", totalAtk)
			fmt.Printf(" atkp: %.4f", atkp)
			fmt.Printf(" cc: %.4f", cr)
			fmt.Printf(" cd: %.4f", cd)
			fmt.Printf(" dmg: %.4f", dmgBonus)
			fmt.Printf(" em: %.4f", em)
			fmt.Print("\n")
			fmt.Print("\tModifiers -")
			fmt.Printf(" talent: %.4f", ab.Talent)
			fmt.Printf(" def adj: %.4f", defAdj)
			fmt.Printf(" res adj: %.4f", resAdj)
			fmt.Printf(" vape/melt: %.4f", vmMult)
			fmt.Print("\n")
		}

		normalDmg := totalAtk * (1 + dmgBonus) * ab.Talent * defAdj * resAdj * vmMult
		critDmg := normalDmg * (1 + cd)
		avgDmg := normalDmg * (1 + (cr * cd))

		r = append(r, DmgResult{
			Normal: normalDmg,
			Avg:    avgDmg,
			Crit:   critDmg,
		})

	}

	return r
}
