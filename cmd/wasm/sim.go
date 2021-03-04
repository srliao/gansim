package main

import (
	"math"
	"math/rand"
	"time"
)

func dmgsim(n, b, w int64, p profile) (start int64, hist []float64, min, max, mean, sd float64) {
	var progress, sum, ss float64
	var data []float64
	min = math.MaxFloat64
	max = -1
	count := n

	resp := make(chan result, n)
	req := make(chan bool)
	done := make(chan bool)
	for i := 0; i < int(w); i++ {
		go dworker(p, resp, req, done)
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

		}
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

func dworker(p profile, resp chan result, req chan bool, done chan bool) {
	seed := time.Now().UnixNano()
	rand := rand.New(rand.NewSource(seed))
	// s.Log.Debugw("worker started", "seed", seed)
	for {
		select {
		case <-req:
			var s artifactSet
			s.Set = make(map[slot]artifact)
			//generate a set of artifacts
			s.Set[Flower] = randArtifact(Flower, HP, p.ArtifactLvl, rand)
			s.Set[Feather] = randArtifact(Feather, ATK, p.ArtifactLvl, rand)
			s.Set[Sands] = randArtifact(Sands, p.ArtifactMainStats[Sands], p.ArtifactLvl, rand)
			s.Set[Goblet] = randArtifact(Goblet, p.ArtifactMainStats[Goblet], p.ArtifactLvl, rand)
			s.Set[Circlet] = randArtifact(Circlet, p.ArtifactMainStats[Circlet], p.ArtifactLvl, rand)

			//calculate dmg
			r := calc(p, s, false)

			var out result
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

func randArtifact(s slot, main statType, lvl int64, rand *rand.Rand) artifact {

	var r artifact

	// r.Slot = s

	// r.MainStat.Type = main
	// r.MainStat.Value = s.MainStat[main][lvl]

	// r.Level = lvl

	// //how many substats
	// p := rand.Float64()
	// var lines = 3
	// if p <= s.FullSubProb {
	// 	lines = 4
	// }

	// //if artifact lvl is less than 4 AND lines =3, then we only want to roll 3 substats
	// n := 4
	// if lvl < 4 && lines < 4 {
	// 	n = 3
	// }

	// //roll initial substats
	// if _, ok := s.SubProb[slot][main]; !ok {
	// 	log.Panicf("main stat %v not found in substat probability map", main)
	// }
	// prb := make(map[StatType]float64)
	// for _, v := range s.SubProb[slot][main] {
	// 	prb[v.Type] = v.Prob
	// }
	// keys := make([]string, len(prb))
	// for k := range prb {
	// 	keys = append(keys, string(k))
	// }
	// sort.Strings(keys)

	// for i := 0; i < n; i++ {
	// 	var sumWeights float64
	// 	for _, v := range prb {
	// 		sumWeights += v
	// 	}
	// 	found := ""
	// 	//pick a number between 0 and sumweights
	// 	pick := rand.Float64() * sumWeights
	// 	for _, k := range keys {
	// 		v := prb[StatType(k)]
	// 		if pick < v && found == "" {
	// 			found = k
	// 		}
	// 		pick -= v
	// 	}
	// 	if found == "" {
	// 		log.Panic("unexpected - no random stat generated")
	// 	}
	// 	t := StatType(found)
	// 	//set prb for this stat to 0 for next iteration
	// 	prb[t] = 0

	// 	tier := rand.Intn(4)
	// 	val := s.SubTier[tier][t]
	// 	r.Substat = append(r.Substat, Stat{
	// 		Type:  t,
	// 		Value: val,
	// 	})
	// }

	// //check how many upgrades to do
	// up := lvl / 4

	// //if we started w 3 lines, then subtract one from # of upgrades
	// if lines == 3 {
	// 	up--
	// }

	// if len(r.Substat) != 4 {
	// 	s.Log.Debugw("invalid artifact, less than 4 lines", "a", r)
	// 	log.Panic("invalid artifact")
	// }

	// //do more rolls, +4/+8/+12/+16/+20
	// for i := 0; i < int(up); i++ {
	// 	pick := rand.Intn(4)
	// 	tier := rand.Intn(4)
	// 	r.Substat[pick].Value += s.SubTier[tier][r.Substat[pick].Type]
	// }

	return r

}
