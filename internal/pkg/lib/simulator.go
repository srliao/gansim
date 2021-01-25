package lib

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//Simulator runs one set of simulation
type Simulator struct {
	MainStat    map[StatType][]float64
	MainProb    map[Slot][]StatProb
	SubTier     []map[StatType]float64
	SubProb     map[Slot]map[StatType][]StatProb //probility of sub stat given main stat
	FullSubProb float64                          //probability of getting 4 lines on an artifact
	Log         *zap.SugaredLogger
	showDebug   bool
}

//NewSimulator creates a new sim
func NewSimulator(
	mainStat map[StatType][]float64,
	mainProb map[Slot][]StatProb,
	subTier []map[StatType]float64,
	subProb map[Slot]map[StatType][]StatProb,
	showDebug bool,
	cfg ...func(*Simulator) error,
) (*Simulator, error) {

	s := &Simulator{
		MainStat:    mainStat,
		MainProb:    mainProb,
		SubTier:     subTier,
		SubProb:     subProb,
		FullSubProb: 207.0 / 932.0,
		showDebug:   showDebug,
	}

	//custom configs
	for _, f := range cfg {
		err := f(s)
		if err != nil {
			return nil, err
		}
	}

	//setup logs
	if s.Log == nil {

		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		if !showDebug {
			config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		}

		logger, err := config.Build()
		if err != nil {
			return nil, err
		}
		sugar := logger.Sugar()
		// sugar.Debugw("logger initiated")

		s.Log = sugar
	}

	//normalize all the probabilities

	for k, x := range s.MainProb {
		var sum float64
		for _, v := range x {
			sum += v.Prob
		}
		for i, v := range x {
			s.MainProb[k][i].Prob = v.Prob / sum
		}
	}

	for slot, x := range s.SubProb {
		for k, y := range x {
			var sum float64
			for _, v := range y {
				sum += v.Prob
			}
			for i, v := range y {
				sum += v.Prob
				s.SubProb[slot][k][i].Prob = v.Prob / sum
			}
		}
	}

	// s.Log.Debugw("prob", "sub prob", s.SubProb)
	// s.Log.Debugw("prob", "main prob", s.MainProb)

	return s, nil
}

//SimDmgDist n = number of sim, b = bin size, w number of worker
func (s *Simulator) SimDmgDist(n, b, w int64, p Profile) (start int64, hist []float64, min, max, mean, sd float64) {
	//calculate the damage distribution
	s.Log.Debugw("starting dmg sim", "n", n, "b", b, "w", w)

	var progress, sum, ss float64
	var data []float64
	min = math.MaxFloat64
	max = -1
	count := n

	resp := make(chan DmgResult, n)
	req := make(chan bool)
	done := make(chan bool)
	for i := 0; i < int(w); i++ {
		go s.workerD(p, resp, req, done)
	}

	//use a go routine to send out a job whenever a worker is done
	go func() {
		var wip int64
		for wip < n {
			//try sending a job to req chan while wip < cfg.NumSim
			req <- true
			wip++
		}
	}()

	if !s.showDebug {
		fmt.Print("\tProgress: 0%")
	}
	for count > 0 {
		//process results received
		r := <-resp
		count--
		val := r.Avg

		//add the avg, rest doesn't really make sense
		data = append(data, val)
		sum += val
		if val < min {
			min = val
		}
		if val > max {
			max = val
		}

		if (1 - float64(count)/float64(n)) > (progress + 0.1) {
			progress = (1 - float64(count)/float64(n))
			if !s.showDebug {
				fmt.Printf("...%.0f%%", 100*progress)
			}
		}
	}
	if !s.showDebug {
		fmt.Print("...100%\n")
	}

	close(done)

	mean = sum / float64(n)
	start = int64(min/float64(b)) * b
	binMax := (int64(max/float64(b)) + 1.0) * b
	numBin := ((binMax - start) / b) + 1

	hist = make([]float64, numBin)

	for _, v := range data {
		ss += (v - mean) * (v - mean)
		steps := int64((v - float64(start)) / float64(b))
		hist[steps]++
	}

	sd = math.Sqrt(ss / float64(n))

	//bin the results
	return

}

//SimArtifactFarm n = number of sim, b = bin size, w number of worker, d damage to exceed, p Profile
func (s *Simulator) SimArtifactFarm(n, b, w int64, d float64, p Profile) (start int64, hist []float64, min, max int64, mean, sd float64) {
	s.Log.Debugw("starting artifact farm sim", "n", n, "b", b, "w", w)

	var progress, ss float64
	var sum int64
	var data []int64
	min = math.MaxInt64
	max = -1
	count := n

	resp := make(chan int64, n)
	req := make(chan bool)
	done := make(chan bool)
	for i := 0; i < int(w); i++ {
		go s.workerA(p, d, resp, req, done)
	}

	//use a go routine to send out a job whenever a worker is done
	go func() {
		var wip int64
		for wip < n {
			//try sending a job to req chan while wip < cfg.NumSim
			req <- true
			wip++
		}
	}()

	if !s.showDebug {
		fmt.Print("\tProgress: 0%")
	}
	for count > 0 {
		//process results received
		r := <-resp
		count--

		//add the avg, rest doesn't really make sense
		data = append(data, r)
		sum += r
		if r < min {
			min = r
		}
		if r > max {
			max = r
		}

		if (1 - float64(count)/float64(n)) > (progress + 0.1) {
			progress = (1 - float64(count)/float64(n))
			if !s.showDebug {
				fmt.Printf("...%.0f%%", 100*progress)
			}
		}
	}
	if !s.showDebug {
		fmt.Print("...100%\n")
	}

	close(done)

	mean = float64(sum) / float64(n)
	binMin := int64(min/b) * b
	binMax := (int64(max/b) + 1) * b
	numBin := (binMax-binMin)/b + 1

	hist = make([]float64, numBin)

	start = binMin

	for _, v := range data {
		ss += (float64(v) - mean) * (float64(v) - mean)
		steps := (v - binMin) / b
		hist[steps]++
	}

	sd = math.Sqrt(ss / float64(n))

	return
}

func (s *Simulator) workerD(p Profile, resp chan DmgResult, req chan bool, done chan bool) {
	seed := time.Now().UnixNano()
	rand := rand.New(rand.NewSource(seed))
	// s.Log.Debugw("worker started", "seed", seed)
	for {
		select {
		case <-req:
			//generate a set of artifacts
			set := make(map[Slot]Artifact)
			set[Flower] = s.RandArtifact(Flower, HP, p.Artifacts.Level, rand)
			set[Feather] = s.RandArtifact(Feather, ATK, p.Artifacts.Level, rand)
			set[Sands] = s.RandArtifact(Sands, p.Artifacts.TargetMainStat[Sands], p.Artifacts.Level, rand)
			set[Goblet] = s.RandArtifact(Goblet, p.Artifacts.TargetMainStat[Goblet], p.Artifacts.Level, rand)
			set[Circlet] = s.RandArtifact(Circlet, p.Artifacts.TargetMainStat[Circlet], p.Artifacts.Level, rand)

			//calculate dmg
			r := Calc(p, set, false)

			var out DmgResult
			for _, v := range r {
				out.Normal += v.Normal
				out.Avg += v.Avg
				out.Crit += v.Crit
			}

			resp <- out
		case <-done:
			return
		}

	}
}

func (s *Simulator) workerA(p Profile, d float64, resp chan int64, req chan bool, done chan bool) {
	seed := time.Now().UnixNano()
	rand := rand.New(rand.NewSource(seed))
	// s.Log.Debugw("worker started", "seed", seed)
	for {
		select {
		case <-req:
			var count int64
			bag := make(map[Slot]Artifact)
			max := -1.0
			/**

			- roll random slot; total++
			  - roll 50/50 if on set
			    - if not on set && is not goblet, discard
				- if not on set && is goblet
				  - roll random main stat; if main stat = 1/5 * ele %, keep; else discard
				- if on set
				  - if not feather/flower, roll random main stat
					- if main stat is atk%, 1/5 * ele %, crit chance, or crit dmg => keep
					- else discard
			  - roll substat if kept
				- if # of cc/cd/atk%/atk < 2, discard
				- else upgrade to +20
			  - calc dmg with kept artifact, if > current, keep new, discard old
			  - if dmg > threshold, stop

			  **/
			//roll random slot
		NEXTTRY:
			for max < d {
				if s.showDebug {
					fmt.Println("------------------------")
				}

				count++
				//panic
				if count > 100000000 {
					log.Fatal("damage not reached after 100mil tries")
				}
				var next Artifact
				//roll 50/50 on set
				onSet := rand.Intn(2) == 0
				//roll random slot
				rs := rand.Intn(5)

				switch rs {
				case 0:
					next.Type = Flower
				case 1:
					next.Type = Feather
				case 2:
					next.Type = Sands
				case 3:
					next.Type = Goblet
				default:
					next.Type = Circlet
				}
				s.Log.Debugw("rand art", "worker", seed, "onSet", onSet, "slot", next.Type, "count", count)
				//if not on set and not goblet, discard
				if !onSet && next.Type != Goblet {
					s.Log.Debugw("rand art", "worker", seed, "discarding offset (not goblet)", next.Type)
					continue NEXTTRY
				}
				//roll random main stat
				rm := rand.Float64()
				prob, ok := s.MainProb[next.Type]
				if !ok {
					log.Fatalf("probability for %v not found", next.Type)
				}
				sum := 0.0
				found := -1
				for i, v := range prob {
					sum += v.Prob
					if rm <= sum && found == -1 {
						found = i
					}
				}
				if found == -1 {
					log.Fatalf("unexpected err generating %v main stat", next.Type)
				}
				ms := s.MainProb[next.Type][found].Type
				s.Log.Debugw("rand art", "worker", seed, "ms", ms)
				//if main stat == ele%, discard 1/6
				if ms == EleP {
					er := rand.Intn(6)
					if er != 0 {
						s.Log.Debugw("rand art", "worker", seed, "discarding 1/6 eleP", er)
						continue NEXTTRY
					}
				}
				//at this point if it's a goblet and ms == eleP then it's the right element for sure
				if !onSet && next.Type == Goblet {
					//if not ele %, discard
					if ms != EleP {
						s.Log.Debugw("rand art", "worker", seed, "discarding offset goblet (not elep)", ms)
						continue NEXTTRY
					}
				}
				//at this point only onSet pieces, or offset correct EleP

				//if not feather or flower, check main stat
				if next.Type != Feather && next.Type != Flower {
					//if main stat is atk%, ele %, crit chance, or crit dmg => keep, else discard
					if ms != ATKP && ms != EleP && ms != CR && ms != CD {
						s.Log.Debugw("rand art", "worker", seed, "discarding non feather/flower", ms)
						continue NEXTTRY
					}
				}

				//roll substat
				next = s.RandArtifact(next.Type, ms, p.Artifacts.Level, rand)

				//if # of cc/cd/atk%/atk < 2, discard
				goodSub := 0
				mustHave := 0
				for _, v := range next.Substat {
					switch v.Type {
					case ATKP:
						goodSub++
						mustHave++
					case CR:
						goodSub++
						mustHave++
					case CD:
						goodSub++
						mustHave++
					case ATK:
						goodSub++
					case ER:
						goodSub++
					}
				}

				s.Log.Debugw("rand art", "worker", seed, "a", next.pretty(), "good", goodSub, "must", mustHave)

				//if sub stat sucks then discard
				if goodSub < 2 || mustHave < 1 {
					s.Log.Debugw("rand art", "worker", seed, "discarding bad stats", fmt.Sprintf("%v - %v", goodSub, mustHave))
					continue NEXTTRY
				}

				var dd float64
				var temp []DmgResult
				//calculate damage full set or not (it'll just be really low with 1 item)

				nextSet := make(map[Slot]Artifact)
				for k, v := range bag {
					nextSet[k] = v
				}
				nextSet[next.Type] = next

				temp = Calc(p, nextSet, false)
				for _, v := range temp {
					dd += v.Avg
				}

				s.Log.Debugw("rand art", "worker", seed, "max", max, "next dmg", dd)
				s.Log.Debugw("rand art", "worker", seed, "current", bag[next.Type].pretty())
				s.Log.Debugw("rand art", "worker", seed, "next", next.pretty())
				s.Log.Debugw("rand art", "worker", seed, "bag", prettySet(bag))
				s.Log.Debugw("rand art", "worker", seed, "next set", prettySet(nextSet))

				if dd > max {
					max = dd
					bag[next.Type] = next
				}
			}

			fmt.Printf("---- Completed in %v, dmg = %v ----\n", count, max)
			fmt.Println(prettySetM(bag))

			artifactStats := make(map[StatType]float64)

			for _, a := range bag {
				artifactStats[a.MainStat.Type] += a.MainStat.Value

				for _, v := range a.Substat {
					artifactStats[v.Type] += v.Value
				}

			}

			s.Log.Debugw("found max", "dmg", max, "count", count, "goblet", bag[Goblet].MainStat.Type, "circlet", bag[Circlet].MainStat.Type, "total stats", artifactStats)

			resp <- count
		case <-done:
			return
		}

	}
}

//RandArtifact generates one random artifact with specified main stat
func (s *Simulator) RandArtifact(slot Slot, main StatType, lvl int64, rand *rand.Rand) Artifact {

	var r Artifact

	r.Type = slot

	r.MainStat.Type = main
	r.MainStat.Value = s.MainStat[main][lvl]

	r.Level = lvl

	//how many substats
	var sum, nextProbSum float64

	var found bool
	p := rand.Float64()
	var lines = 3
	if p <= s.FullSubProb {
		lines = 4
	}
	//roll initial substats
	prb, ok := s.SubProb[slot][main]
	if !ok {
		log.Panicf("main stat %v not found in substat probability map", main)
	}
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
		p = rand.Float64()
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
				tier := rand.Intn(4)
				val := s.SubTier[tier][v.Type]
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

	if len(r.Substat) != 4 {
		s.Log.Debugw("invalid artifact, less than 4 lines", "a", r)
		log.Panic("invalid artifact")
	}

	//do more rolls, +4/+8/+12/+16/+20
	for i := 0; i < int(up); i++ {
		pick := rand.Intn(4)
		tier := rand.Intn(4)
		r.Substat[pick].Value += s.SubTier[tier][r.Substat[pick].Type]
	}

	return r

}
