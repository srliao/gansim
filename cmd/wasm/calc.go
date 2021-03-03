package main

import (
	"fmt"
	"log"
	"math"
)

type result struct {
	Normal float64 `json:"normal"`
	Avg    float64 `json:"avg"`
	Crit   float64 `json:"crit"`
}

func calc(p profile, s artifactSet, showDebug bool) []result {
	var r []result

	//add up artifact stats
	artifactStats := make(map[statType]float64)

	for _, a := range s.Set {
		artifactStats[a.MainStat.Type] += a.MainStat.Value
		for _, v := range a.Substat {
			artifactStats[v.Type] += v.Value
		}
	}

	charMods := make(map[modType]float64)
	for _, a := range p.Character.Mods {
		for _, v := range a.List {
			charMods[a.Type] = charMods[a.Type] + v.Value
		}
	}

	weaponMods := make(map[modType]float64)
	for _, a := range p.Weapon.Mods {
		for _, v := range a.List {
			weaponMods[a.Type] = weaponMods[a.Type] + v.Value
		}
	}

	baseAtk := p.Character.BaseAtk + p.Weapon.BaseAtk

	if showDebug {
		fmt.Print("--------------------------\n")
		fmt.Printf("\tCharacter + weapon base atk: %.2f\n", baseAtk)
		fmt.Printf("\t%v\n", artifactStats)
	}

	for _, ab := range p.Abilities {
		var totalAtk, atk, atkp, cr, cd, elep, phyp, em, dmgBonus, reactionBonus, defShred, defAdj, resAdj, res float64
		totalAtk = baseAtk

		atk += artifactStats[ATK]
		atkp += artifactStats[ATKP]
		cr += artifactStats[CR]
		cd += artifactStats[CD]
		elep += artifactStats[EleP]
		em += artifactStats[EM]
		phyp += artifactStats[PhyP]

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

		//map abil mods
		abilMods := make(map[modType]float64)
		for _, a := range ab.Mods {
			for _, v := range a.List {
				abilMods[a.Type] = abilMods[a.Type] + v.Value
			}
		}

		//add up all the mods
		atk += charMods[ATKPMod] + weaponMods[ATKPMod] + abilMods[ATKPMod]
		atkp += charMods[ATKPMod] + weaponMods[ATKPMod] + abilMods[ATKPMod]
		cr += charMods[CRMod] + weaponMods[CRMod] + abilMods[CRMod]
		cd += charMods[CDMod] + weaponMods[CDMod] + abilMods[CDMod]
		elep += charMods[ElePMod] + weaponMods[ElePMod] + abilMods[ElePMod]
		em += charMods[EMMod] + weaponMods[EMMod] + abilMods[EMMod]
		phyp += charMods[PhyPMod] + weaponMods[PhyPMod] + abilMods[PhyPMod]

		//add up dmg mods
		//only if abil is element
		if ab.IsPhys {
			dmgBonus += phyp
		} else {
			dmgBonus += elep
		}

		totalAtk = totalAtk*(1+atkp) + atk

		//cap cc at 1
		if cr > 1 {
			cr = 1
		}
		if cr < 0 {
			log.Println("WARNING, CRIT RATE < 0")
			cr = 0
		}
		if cd < 0 {
			cd = 0
			log.Println("WARNING, CRIT DMG < 0")
		}

		//add up def shreds
		defShred += charMods[DefMod] + weaponMods[DefMod] + abilMods[DefMod]
		//calculate def adjustment
		defAdj = float64(100+p.Character.Level) / (float64(100+p.Character.Level) + float64(100+p.Enemy.Level)*(1-defShred))

		//calculate res adjustment
		res += charMods[ResMod] + weaponMods[ResMod] + abilMods[ResMod]

		if res < 0 {
			resAdj = 1 - (res / 2)
		} else if res < 0.75 {
			resAdj = 1 - res
		} else {
			resAdj = 1 / (4*res + 1)
		}

		//calculate reaction bonus
		reactionBonus += charMods[ReactMod] + weaponMods[ReactMod] + abilMods[ReactMod]

		vmMult := 1.0
		if ab.IsVapeMelt {
			vmMult = ab.VapeMeltMul*(1+0.00189266831*em*math.Exp(-0.000505*em)) + reactionBonus
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
			fmt.Printf(" talent: %.4f", ab.Multiplier)
			fmt.Printf(" def adj: %.4f", defAdj)
			fmt.Printf(" res adj: %.4f", resAdj)
			fmt.Printf(" vape/melt: %.4f", vmMult)
			fmt.Print("\n")
		}

		normalDmg := totalAtk * (1 + dmgBonus) * ab.Multiplier * defAdj * resAdj * vmMult
		critDmg := normalDmg * (1 + cd)
		avgDmg := normalDmg * (1 + (cr * cd))

		if showDebug {
			fmt.Print("--------------------------\n")
			fmt.Print("\tResults -")
			fmt.Printf(" normal: %.4f", normalDmg)
			fmt.Printf(" avg: %.4f", avgDmg)
			fmt.Printf(" crit: %.4f", critDmg)
			fmt.Print("\n")
		}

		r = append(r, result{
			Normal: normalDmg,
			Avg:    avgDmg,
			Crit:   critDmg,
		})

	}

	return r
}
