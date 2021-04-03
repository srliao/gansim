package upgrade

import (
	"log"
	"math/rand"

	"go.uber.org/zap"
)

type Simulator struct {
	*zap.SugaredLogger
}

func New(p Profile) (*Simulator, error) {
	return nil, nil
}

func (a *Simulator) rand(rand *rand.Rand, main StatType, lvl int) ([]StatType, []float64, float64) {
	stats := make([]float64, 4)

	//how many substats
	var lines = 3
	if rand.Float64() <= 306.0/1415.0 {
		lines = 4
	}
	//if artifact lvl is less than 4 AND lines =3, then we only want to roll 3 substats
	n := 4
	if lvl < 4 && lines < 4 {
		n = 3
	}
	//make a copy of prob
	prb := make([]float64, len(weights))

	for i, v := range weights {
		if i == int(main) {
			v = 0
		}
		prb[i] = v
	}

	subType := []StatType{-1, -1, -1, -1}

	for i := 0; i < n; i++ {
		var sumWeights float64
		for _, v := range prb {
			sumWeights += v
		}
		found := -1
		//pick a number between 0 and sumweights
		pick := rand.Float64() * sumWeights
		for i, v := range prb {
			if pick < v && found == -1 {
				found = i
			}
			pick -= v
		}
		if found == -1 {
			log.Println("sum weights ", sumWeights)
			log.Println("prb ", prb)
			log.Panic("unexpected - no random stat generated")
		}
		// log.Println("found at ", found)
		subType[i] = StatType(found)
		//set prb for this stat to 0 for next iteration
		prb[found] = 0

		tier := rand.Intn(4)

		stats[i] += tiers[found][tier]
	}

	//check how many upgrades to do
	up := lvl / 4

	//if we started w 3 lines, then subtract one from # of upgrades
	if lines == 3 {
		up--
	}

	//do more rolls, +4/+8/+12/+16/+20
	for i := 0; i < int(up); i++ {
		pick := rand.Intn(4)
		tier := rand.Intn(4)

		stats[pick] += tiers[subType[pick]][tier]
	}

	//figure out main stat based on lvl
	m := mainStat[main][lvl]

	return subType, stats, m
}

//Slot identifies the artifact slot
type Slot int

//Types of artifact slots
const (
	Flower Slot = iota
	Feather
	Sands
	Goblet
	Circlet
)

func (s Slot) String() string {
	return [...]string{
		"Flower",
		"Feather",
		"Sands",
		"Goblet",
		"circlet",
	}[s]
}

type StatType int

//stat types
const (
	HP StatType = iota
	ATK
	DEF
	HPP
	ATKP
	DEFP
	ER
	EM
	CR
	CD
	Heal
	EleP
	PhyP
)

func (s StatType) String() string {
	return StatTypeString[s]
}

var StatTypeString = [...]string{
	"HP",
	"ATK",
	"DEF",
	"HP%",
	"ATK%",
	"DEF%",
	"ER",
	"EM",
	"CR",
	"CD",
	"Heal",
	"Ele%",
	"Phys%",
}

var weights = []float64{
	150,
	150,
	150,
	100,
	100,
	100,
	100,
	100,
	75,
	75,
}

var tiers = [][]float64{
	{209, 239, 269, 299},         //hp
	{14, 16, 18, 19},             //atk
	{16, 19, 21, 23},             //def
	{0.041, 0.047, 0.053, 0.058}, //hp%
	{0.041, 0.047, 0.053, 0.058}, //atk%
	{0.051, 0.058, 0.066, 0.073}, //def%
	{0.045, 0.052, 0.058, 0.065}, //er
	{16, 19, 21, 23},             //em
	{0.027, 0.031, 0.035, 0.039}, //cr
	{0.054, 0.062, 0.07, 0.078},  //cd
}
