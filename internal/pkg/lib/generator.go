package lib

import (
	"math/rand"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	ShowDebug       bool
	Log             *zap.SugaredLogger
}

//NewGenerator creates a new artifact generator
func NewGenerator(
	seed int64,
	mainStatLvls map[StatType][]float64,
	subTier []map[StatType]float64,
	mainProb map[StatType]float64,
	subProb map[StatType][]StatProb,
	cfg ...func(*Generator) error,
) (*Generator, error) {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	g := &Generator{
		rand:            r,
		mainStatLvls:    mainStatLvls,
		subTier:         subTier,
		mainProb:        mainProb,
		subProb:         subProb,
		fullSubstatProb: 207.0 / 932.0,
	}

	//custom configs
	for _, f := range cfg {
		err := f(g)
		if err != nil {
			return nil, err
		}
	}

	//setup logs
	if g.Log == nil {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, err := config.Build()
		if err != nil {
			return nil, err
		}
		sugar := logger.Sugar()
		sugar.Debugw("logger initiated")

		g.Log = sugar
	}

	return g, nil
}

//Rand generates one random artifact
func (g *Generator) Rand(lvl int64) Artifact {
	// var r Artifact
	p := g.rand.Float64()
	sum := 0.0
	var main StatType
	for k, v := range g.mainProb {
		sum += v
		if p <= sum {
			main = k
			break
		}
	}
	r := g.RandWithMain(main, lvl)
	return r
}

//RandWithMain generates one random artifact with specified main stat
func (g *Generator) RandWithMain(main StatType, lvl int64) Artifact {

	var r Artifact

	r.MainStat.Type = main
	r.MainStat.Value = g.mainStatLvls[main][lvl]

	//how many substats
	var sum, nextProbSum float64

	var found bool
	p := g.rand.Float64()
	var lines = 3
	if p <= g.fullSubstatProb {
		lines = 4
	}
	//roll initial substats
	prb := g.subProb[main]
	var next []StatProb
	for _, v := range prb {
		nextProbSum += v.Prob
		next = append(next, v)
	}

	//if artifact lvl is less than 4 AND lines =3, then we only want to roll 3 substats
	n := 4
	if lvl < 4 && lines < 4 {
		n = 3
	}

	for i := 0; i < n; i++ {
		var current []StatProb
		var check float64
		for _, v := range next {
			current = append(current, StatProb{Type: v.Type, Prob: v.Prob / nextProbSum})
			check += v.Prob / nextProbSum
		}
		if g.ShowDebug {
			g.Log.Debugw("generating substat", "current prob", current, "count", len(current), "total prob", check)
		}
		p = g.rand.Float64()
		//reset next
		next = []StatProb{}
		nextProbSum = 0
		sum = 0
		found = false
		for _, v := range current {
			sum += v.Prob
			if p <= sum && !found {
				//this is the one!
				//roll 1 to 4 for tier
				//ASSUMPTION = equal weight for each tier
				tier := g.rand.Intn(4)
				val := g.subTier[tier][v.Type]
				r.Substat = append(r.Substat, Stat{
					Type:  v.Type,
					Value: val,
				})
				found = true
			} else {
				//add this one so it's available for next roll
				next = append(next, v)
				nextProbSum += v.Prob
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
		pick := g.rand.Intn(4)
		tier := g.rand.Intn(4)
		r.Substat[pick].Value += g.subTier[tier][r.Substat[pick].Type]
	}

	return r

}

//RandSet generates one set of random artifact
func (g *Generator) RandSet() {

}

//RandSetWithMain generates set of one random artifact with specified main stat
func (g *Generator) RandSetWithMain() {

}
