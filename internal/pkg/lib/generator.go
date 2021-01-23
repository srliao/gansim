package lib

import (
	"fmt"
	"math/rand"
	"time"
)

//Stat represents one stat
type Stat struct {
	Type  StatType
	Value float64
}

//StatProb represents probability of getting a stat
type StatProb struct {
	Type StatType
	Prob float64
}

//Generator for generating random artifacts
type Generator struct {
	rand            *rand.Rand
	mainStatLvls    map[StatType][]float64
	subTier         []map[StatType]float64
	subProb         map[StatType][]StatProb //probility of sub stat given main stat
	mainProb        map[StatType]float64
	fullSubstatProb float64 //probability of getting 4 lines on an artifact
}

//NewGenerator creates a new artifact generator
func NewGenerator(
	seed int64,
	mainStatLvls map[StatType][]float64,
	subTier []map[StatType]float64,
	mainProb map[StatType]float64,
	subProb map[StatType][]StatProb,
) *Generator {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	generator := Generator{
		rand:            r,
		mainStatLvls:    mainStatLvls,
		subTier:         subTier,
		mainProb:        mainProb,
		subProb:         subProb,
		fullSubstatProb: 207.0 / 932.0,
	}

	return &generator
}

//Rand generates one random artifact
func (g *Generator) Rand(lvl int64) {
	// var r Artifact

}

//RandWithMain generates one random artifact with specified main stat
func (g *Generator) RandWithMain(main StatType, lvl int64) {

	var r Artifact

	r.MainStat.Type = main
	r.MainStat.Value = g.mainStatLvls[main][lvl]

	//how many substats
	var prb4stat, sum, nextProbSum float64
	var next []StatProb
	var found bool
	p := g.rand.Float64()
	var lines = 3
	if p <= g.fullSubstatProb {
		lines = 4
	}
	//roll initial substats
	//line 1
	//make sure prob adds to 1 first

	for _, v := range prb {
		nextProbSum += v.p
		next = append(next, v)
	}

	//if artifact lvl is less than 4 AND lines =3, then we only want to roll 3 substats
	n := 4
	if lvl < 4 && lines < 4 {
		n = 3
	}

	for i := 0; i < n; i++ {
		var current []statPrb
		var check float64
		for _, v := range next {
			current = append(current, statPrb{t: v.t, p: v.p / nextProbSum})
			check += v.p / nextProbSum
		}
		if showDebug {
			fmt.Println("current probabilities: ", current)
			fmt.Println("sub stat count: ", len(current))
			fmt.Println("current prob total: ", check)
		}
		p = generator.Float64()
		next = []statPrb{}
		nextProbSum = 0
		sum = 0
		found = false
		for _, v := range current {
			sum += v.p
			if p <= sum && !found {
				//this is the one!
				//roll 1 to 4 for tier
				//ASSUMPTION = equal weight for each tier
				tier := generator.Intn(4)
				val := cfg.SubstatTier[tier][v.t]
				r.Sub = append(r.Sub, stat{
					t: v.t,
					s: val,
				})
				found = true
			} else {
				//add this one so it's available for next roll
				next = append(next, v)
				nextProbSum += v.p
			}
		}
	}

	//check how many upgrades to do
	up := lvl / 4

	//if we started w 3 lines, then subtract one from # of upgrades
	if lines == 3 {
		up--
	}

	//do more rolls, +4/+8/+12/+16/+20
	for i := 0; i < int(up); i++ {
		pick := generator.Intn(4)
		tier := generator.Intn(4)
		r.Sub[pick].s += cfg.SubstatTier[tier][r.Sub[pick].t]
	}

	return r

}

//RandSet generates one set of random artifact
func (g *Generator) RandSet() {

}

//RandSetWithMain generates set of one random artifact with specified main stat
func (g *Generator) RandSetWithMain() {

}
