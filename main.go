package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"gopkg.in/yaml.v2"
)

const showDebug = false
const testCurrent = false

type statTypes string

const (
	//stat types
	sDEFP statTypes = "DEF%"
	sDEF  statTypes = "DEF"
	sHP   statTypes = "HP"
	sHPP  statTypes = "HP%"
	sATK  statTypes = "ATK"
	sATKP statTypes = "ATK%"
	sER   statTypes = "ER"
	sEM   statTypes = "EM"
	sCC   statTypes = "CR"
	sCD   statTypes = "CD"
	sHeal statTypes = "Heal"
	sEleP statTypes = "Ele%"
	sPhyP statTypes = "Phys%"
)

type config struct {
	Profiles        []string `yaml:"Profiles"`
	GraphOutput     string   `yaml:"GraphOutput"`
	NumSim          int64    `yaml:"NumSim"`
	BinSize         int64    `yaml:"BinSize"`
	WriteCSV        bool     `yaml:"WriteCSV"`
	DamageType      string   `yaml:"DamageType"`
	NumWorker       int64    `yaml:"NumWorker"`
	MainStatFile    string   `yaml:"MainStatFile"`
	SubstatTierFile string   `yaml:"SubstatTierFile"`
	MainStatScaling map[statTypes][]float64
	SubstatTier     []map[statTypes]float64
}

type profile struct {
	Output string `yaml:"Output"`
	Label  string `yaml:"Label"`
	//character info
	CharLevel     float64 `yaml:"CharacterLevel"`
	CharBaseAtk   float64 `yaml:"CharacterBaseAtk"`
	WeaponBaseAtk float64 `yaml:"WeaponBaseAtk"`
	EnemyLevel    float64 `yaml:"EnemyLevel"`
	//artifact info
	ArtifactMaxLevel int64     `yaml:"ArtifactMaxLevel"`
	Sands            statTypes `yaml:"Sands"`
	Goblet           statTypes `yaml:"Goblet"`
	Circlet          statTypes `yaml:"Circlet"`
	SubstatFile      string    `yaml:"SubstatFile"`
	SubstatWeights   map[string][]statPrb
	//abilities
	Abilities []struct {
		Talent      float64   `yaml:"Talent"`
		AtkMod      []float64 `yaml:"AtkMod"`
		EleMod      []float64 `yaml:"EleMod"`
		CCMod       []float64 `yaml:"CCMod"`
		CDMod       []float64 `yaml:"CDMod"`
		DmgMod      []float64 `yaml:"DmgMod"`
		ResistMod   []float64 `yaml:"ResistMod"`
		DefShredMod []float64 `yaml:"DefShredMod"`
	} `yaml:"Abilities"`
}

type stat struct {
	t statTypes
	s float64
}

type statPrb struct {
	t statTypes
	p float64
}

type artifact struct {
	Main stat
	Sub  []stat //4
}

type artifacts struct {
	Flower  artifact
	Feather artifact
	Sands   artifact
	Goblet  artifact
	Circlet artifact
}

func main() {
	// runtime.GOMAXPROCS(12)
	//read config

	var err error

	f, err := os.Open("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// fmt.Println("reading config file")
	var cfg config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(cfg)

	subtier, err := loadSubstatTier(cfg.SubstatTierFile)
	if err != nil {
		log.Fatal(err)
	}

	cfg.SubstatTier = subtier

	msscaling, err := loadMainStatScaling(cfg.MainStatFile)
	if err != nil {
		log.Fatal(err)
	}
	cfg.MainStatScaling = msscaling

	//loop through profiles and run sim for each

	histData := make([][]float64, len(cfg.Profiles))
	labels := make([]string, len(cfg.Profiles))

	for i, ppath := range cfg.Profiles {

		source, err := ioutil.ReadFile(ppath)
		if err != nil {
			log.Fatal(err)
		}

		// log.Printf("reading profile: %v\n", ppath)
		var prf profile

		err = yaml.Unmarshal(source, &prf)
		if err != nil {
			log.Fatal(err)
		}

		prf.SubstatWeights = make(map[string][]statPrb)

		//load substat weights
		sw, err := os.Open(prf.SubstatFile)
		if err != nil {
			log.Fatal(err)
		}
		defer sw.Close()

		reader := csv.NewReader(sw)
		lines, err := reader.ReadAll()
		if err != nil {
			log.Fatal(err)
		}
		var headers []string
		for i, line := range lines {
			if i == 0 {
				//read the headers
				for j := 1; j < len(line); j++ {
					headers = append(headers, line[j])
				}
				// fmt.Printf("headers: %v\n", headers)
			} else {
				//otherwise populate
				t := line[0]
				for j := 1; j < len(line); j++ {
					//parse float
					p, err := strconv.ParseFloat(line[j], 64)
					if err != nil {
						log.Println("err parsing float at line ", i)
						log.Fatal(err)
					}
					prf.SubstatWeights[headers[j-1]] = append(prf.SubstatWeights[headers[j-1]], statPrb{t: statTypes(t), p: p})
				}
			}
		}

		if testCurrent {
			test(prf)
			// return?
		}

		// fmt.Printf("%v\n", prf.Abilities)

		labels[i] = fmt.Sprintf("%v [lvl %v]", prf.Label, prf.ArtifactMaxLevel)

		fmt.Printf("starting simulation for profile: %v, n = %v\n", ppath, cfg.NumSim)
		timeStart := time.Now()
		hist := sim(cfg, prf)
		elapsed := time.Since(timeStart)
		fmt.Printf("Simulation for profile %v took %s\n\n", ppath, elapsed)
		histData[i] = hist

	}

	// fmt.Println("labels: ", labels)

	//sim results page
	items := make([][]opts.LineData, len(cfg.Profiles))
	min := make([]float64, len(cfg.Profiles))
	max := make([]float64, len(cfg.Profiles))
	avg := make([]float64, len(cfg.Profiles))
	ss := make([]float64, len(cfg.Profiles))
	bins := make([][]float64, len(cfg.Profiles))

	//absolute max/min
	var binMax, binMin float64
	binMax = -100000000
	binMin = 100000000

	for i := range min {
		min[i] = 1000000000
		max[i] = -1000000000
	}

	//find the min/max/avg/std
	for i, hist := range histData {
		for _, v := range hist {
			if v < min[i] {
				min[i] = v
			}
			if v > max[i] {
				max[i] = v
			}
			avg[i] += v
		}
		avg[i] = avg[i] / float64(len(hist))
		if binMin > min[i] {
			binMin = min[i]
		}
		if binMax < max[i] {
			binMax = max[i]
		}
	}
	//calculate bin size
	binMin = float64(int64(binMin/float64(cfg.BinSize)) * cfg.BinSize)
	binMax = float64((int64(binMax/float64(cfg.BinSize)) + 1.0) * cfg.BinSize)

	// fmt.Printf("bin min: %v, bin max: %v\n", binMin, binMax)

	numBin := int64((binMax-binMin)/float64(cfg.BinSize)) + 1
	xaxis := make([]float64, numBin)

	//bin the data
	for i, hist := range histData {
		bins[i] = make([]float64, numBin)
		for _, v := range hist {
			ss[i] += (v - avg[i]) * (v - avg[i])
			//find the steps and bin this
			steps := int64((v - float64(binMin)) / float64(cfg.BinSize))
			bins[i][steps]++
		}
	}

	for i, b := range bins {
		for j, v := range b {
			items[i] = append(items[i], opts.LineData{Value: v / float64(cfg.NumSim)})
			xaxis[j] = binMin + float64(j)*float64(cfg.BinSize) + float64(cfg.BinSize/2)
		}
	}

	//add min, max, avg, stddev to label
	for i, v := range labels {
		sd := math.Sqrt(ss[i] / float64(cfg.NumSim))
		labels[i] = fmt.Sprintf("%v (min: %.f max: %.f avg: %.f sd: %.f)", v, min[i], max[i], avg[i], sd)
	}

	page := components.NewPage()
	page.PageTitle = "simulation results"
	lineChart := charts.NewLine()
	lineChart.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: fmt.Sprintf("Probability Density Function (n = %v)", cfg.NumSim),
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Probability",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: fmt.Sprintf("Dmg: %v", cfg.DamageType),
		}),
		// charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithLegendOpts(opts.Legend{Show: true, Right: "0%", Orient: "vertical", Data: labels}),
	)
	lineChart.SetXAxis(xaxis)

	//add items to our chart
	for i, series := range items {
		lineChart.AddSeries(labels[i], series)
	}

	//add all hist data into the charts
	page.AddCharts(
		lineChart,
	)
	graph, err := os.Create(cfg.GraphOutput)
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(graph))

}

func sim(cfg config, p profile) []float64 {

	n := cfg.NumSim
	writeCSV := cfg.WriteCSV
	//generate some random sets
	var avg, min, max, progress float64
	count := n
	min = 10000000 //shouldnt ever get this big...
	var hist []float64

	// wc := 0
	// wcmax := cfg.NumWorker
	resp := make(chan result, cfg.NumSim)

	// fmt.Println("n = ", count)
	//keep sending out workers while still simulations left to do
	fmt.Print("\tProgress: 0%")

	//fire up max number of workers
	req := make(chan bool)
	done := make(chan bool)
	for i := 0; i < int(cfg.NumWorker); i++ {
		go worker(cfg, p, resp, req, done)
	}

	//use a go routine to send out a job whenever a worker is done
	go func() {
		var wip int64
		for wip < n {
			//try sending a job to req chan while wip < n
			req <- true
			wip++
		}
	}()

	for count > 0 {
		//process results received
		r := <-resp
		count--

		val := r.a
		switch cfg.DamageType {
		case "normal":
			val = r.n
		case "crit":
			val = r.c
		}

		//push result to r
		hist = append(hist, val)
		if val < min {
			min = val
		}
		if val > max {
			max = val
		}
		avg += val
		if (1 - float64(count)/float64(n)) > (progress + 0.1) {
			progress = (1 - float64(count)/float64(n))
			fmt.Printf("...%.0f%%", 100*progress)
		}
	}
	fmt.Print("...100%%\n")

	close(done)

	if writeCSV && p.Output != "" {

		avg = avg / float64(n)

		//bin it in 200 increments, starting at min rounded down to nearest 200 up to max rounded up to nearest 200
		var inc int64
		inc = 200
		binMin := int64(min/float64(inc)) * inc
		binMax := (int64(max/float64(inc)) + 1) * inc
		numBin := (binMax - binMin) / inc

		bins := make([]float64, numBin)

		var ss float64

		for _, v := range hist {
			steps := int64((v - float64(binMin)) / float64(inc))
			bins[steps]++
			//calculate the std dev while we're at it
			ss += (v - avg) * (v - avg)
		}

		std := math.Sqrt(ss / float64(n))

		os.Remove(p.Output)
		file, err := os.Create(p.Output)
		if err != nil {
			log.Panicln(err)
		}
		defer file.Close()
		writer := csv.NewWriter(file)
		defer writer.Flush()

		for i, v := range bins {
			start := strconv.FormatInt(binMin+int64(i)*inc, 10)
			end := strconv.FormatInt(binMin+int64(i+1)*inc, 10)

			if err := writer.Write([]string{
				strconv.FormatInt(int64(i), 10),
				start,
				end,
				strconv.FormatInt(int64(v), 10),
			}); err != nil {
				log.Panicln(err)
			}
		}

		if err := writer.Write([]string{
			"avg",
			strconv.FormatFloat(avg, 'f', 6, 64),
		}); err != nil {
			log.Panicln(err)
		}

		if err := writer.Write([]string{
			"std dev",
			strconv.FormatFloat(std, 'f', 6, 64),
		}); err != nil {
			log.Panicln(err)
		}

		fmt.Printf("std dev: %v, avg: %v, min: %v, max: %v\n", std, avg, min, max)
	}

	return hist
}

type result struct {
	n float64
	c float64
	a float64
}

func worker(cfg config, p profile, resp chan result, req chan bool, done chan bool) {
	//create our own rand source
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	for {
		select {
		case <-req:
		case <-done:
			return
		}
		art := genArtifacts(cfg, p, generator)
		if showDebug {
			fmt.Printf("artifacts: %v\n", art)
		}
		//bloom dmg for now
		r := calc(art, p)
		//add up the results
		var total result
		for _, v := range r {
			total.a += v.a
			total.c += v.c
			total.n += v.n
		}
		resp <- total
	}
}

//generate a set of artifacts given the configs
func genArtifacts(cfg config, p profile, generator *rand.Rand) artifacts {
	var r artifacts
	//flower is always hp = 4780
	r.Flower = randArtifact(
		cfg,
		p.ArtifactMaxLevel,
		stat{
			t: sHP,
			s: cfg.MainStatScaling[sHP][p.ArtifactMaxLevel],
		},
		p.SubstatWeights["flower"],
		generator,
	)
	//feather is always flat atk 311
	r.Feather = randArtifact(
		cfg,
		p.ArtifactMaxLevel,
		stat{
			t: sATK,
			s: cfg.MainStatScaling[sATK][p.ArtifactMaxLevel],
		},
		p.SubstatWeights["feather"],
		generator,
	)
	//sands is always % atk 46.6%
	r.Sands = randArtifact(
		cfg,
		p.ArtifactMaxLevel,
		stat{
			t: p.Sands,
			s: cfg.MainStatScaling[p.Sands][p.ArtifactMaxLevel],
		},
		p.SubstatWeights["sands"],
		generator,
	)
	//goblet is always % ele 46.6%
	// r.Goblet = randArtifact(cfg, stat{t: p.Goblet, s: p.GobletStat}, p.SubstatWeights["goblet"], generator)
	r.Goblet = randArtifact(
		cfg,
		p.ArtifactMaxLevel,
		stat{
			t: p.Goblet,
			s: cfg.MainStatScaling[p.Goblet][p.ArtifactMaxLevel],
		},
		p.SubstatWeights["goblet"],
		generator,
	)
	//circlet is always crit dmg 62.20%
	r.Circlet = randArtifact(
		cfg,
		p.ArtifactMaxLevel,
		stat{
			t: p.Circlet,
			s: cfg.MainStatScaling[p.Circlet][p.ArtifactMaxLevel],
		},
		p.SubstatWeights["circlet"],
		generator,
	)
	maxUp := p.ArtifactMaxLevel / 4
	//do some sanity checks on sub stat
	for _, v := range r.Flower.Sub {
		max := cfg.SubstatTier[3][v.t] * float64(maxUp)
		if v.s > max {
			log.Panicf("invalid flower detected, substat %v exceed 5x max tier. %v\n", v.t, r)
		}
		if v.t == r.Flower.Main.t {
			log.Panicf("invalid flower detected, substat %v is same as main stat. %v\n", v.t, r)
		}
	}
	for _, v := range r.Feather.Sub {
		max := cfg.SubstatTier[3][v.t] * float64(maxUp)
		if v.s > max {
			log.Panicf("invalid feather detected, substat %v exceed 5x max tier. %v\n", v.t, r)
		}
		if v.t == r.Feather.Main.t {
			log.Panicf("invalid feather detected, substat %v is same as main stat. %v\n", v.t, r)
		}
	}
	for _, v := range r.Sands.Sub {
		max := cfg.SubstatTier[3][v.t] * float64(maxUp)
		if v.s > max {
			log.Panicf("invalid sands detected, substat %v exceed 5x max tier. %v\n", v.t, r)
		}
		if v.t == r.Sands.Main.t {
			log.Panicf("invalid sands detected, substat %v is same as main stat. %v\n", v.t, r)
		}
	}
	for _, v := range r.Goblet.Sub {
		max := cfg.SubstatTier[3][v.t] * float64(maxUp)
		if v.s > max {
			log.Panicf("invalid goblet detected, substat %v exceed 5x max tier. %v\n", v.t, r)
		}
		if v.t == r.Goblet.Main.t {
			log.Panicf("invalid goblet detected, substat %v is same as main stat. %v\n", v.t, r)
		}
	}
	for _, v := range r.Circlet.Sub {
		max := cfg.SubstatTier[3][v.t] * float64(maxUp)
		if v.s > max {
			log.Panicf("invalid circlet detected, substat %v exceed 5x max tier. %v\n", v.t, r)
		}
		if v.t == r.Circlet.Main.t {
			log.Panicf("invalid circlet detected, substat %v is same as main stat. %v\n", v.t, r)
		}
	}
	return r
}

//randArtifact creates a random artifact given the main stat, assume statPrb adds up to 1
func randArtifact(cfg config, lvl int64, main stat, prb []statPrb, generator *rand.Rand) artifact {
	var r artifact
	r.Main.s = main.s
	r.Main.t = main.t
	//how many substats
	var prb4stat, sum, nextProbSum float64
	var next []statPrb
	var found bool
	prb4stat = 207 / 932
	p := generator.Float64()
	var lines = 3
	if p <= prb4stat {
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

func calc(a artifacts, p profile) []result {

	//artifact substats
	artSubStat := a.substatTotal()

	var r []result

	//loop through each talent
	for _, ab := range p.Abilities {
		//calculate total atk
		var totalAtk, atkp, cc, cd, atk, elep, dmgBonus, defAdj, resAdj, resist float64
		//base atk
		totalAtk = p.CharBaseAtk + p.WeaponBaseAtk
		atk += artSubStat[sATK]
		atkp += artSubStat[sATKP]
		elep += artSubStat[sEleP]
		cc += artSubStat[sCC]
		cd += artSubStat[sCD]

		if showDebug {
			fmt.Printf("\tbase stats - atk: %.4f\n", totalAtk)
			fmt.Printf("\tartifact stats - atk: %.4f, atkp: %.4f, elep: %.4f, cc: %.4f, cd:%.4f\n", atk, atkp, elep, cc, cd)
		}

		//add up def shreds
		var defShred float64
		for _, v := range ab.DefShredMod {
			defShred += v
		}
		//calculate def adjustment
		defAdj = (100 + p.CharLevel) / ((100 + p.CharLevel) + (100+p.EnemyLevel)*(1-defShred))

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
			cc += v
		}
		for _, v := range ab.CDMod {
			cd += v
		}

		//cap cc at 1
		if cc > 1 {
			cc = 1
		}
		if cc < 0 {
			cc = 0
		}
		if cd < 0 {
			cd = 0
		}

		//calculate enemy resistance
		//TODO: not entirely sure this formula is correct??
		//this formula suggests something diff: https://www.reddit.com/r/Genshin_Impact/comments/krg2ic/the_complete_genshin_impact_damage_formula/
		for _, v := range ab.ResistMod {
			resist += v
			// //if v will bring resist to negative, the half the portion that brings it to negative
			// if resist+v < 0 {
			// 	//if resist is already negative
			// 	if resist < 0 {
			// 		//if v is positive?? just don't list it this way
			// 		if v >= 0 {
			// 			resist += v
			// 		} else {
			// 			//half the effect of v
			// 			resist += v / 2
			// 		}
			// 	} else {
			// 		temp := v + resist
			// 		resist = 0
			// 		resist += temp / 2
			// 	}
			// } else {
			// 	resist += v
			// }
		}

		if resist < 0 {
			resAdj = 1 - (resist / 2)
		} else if resist < 0.75 {
			resAdj = 1 - resist
		} else {
			resAdj = 1 / (4*resist + 1)
		}

		if showDebug {
			fmt.Printf("\ttotal attack percent after mod: %.4f\n", atkp)
			fmt.Printf("\ttotal sheet attack: %.4f\n", totalAtk)
			fmt.Printf("\ttotal dmg mod: %.4f\n", dmgBonus)
			fmt.Printf("\ttotal cc: %.4f\n", cc)
			fmt.Printf("\ttotal cd: %.4f\n", cd)
			fmt.Printf("\tdef adj: %.4f\n", defAdj)
			fmt.Printf("\tenemy resist: %.4f\n", resist)
		}

		normalDmg := totalAtk * (1 + dmgBonus) * ab.Talent * defAdj * resAdj
		critDmg := normalDmg * (1 + cd)
		avgDmg := normalDmg * (1 + (cc * cd))

		r = append(r, result{
			n: normalDmg,
			a: avgDmg,
			c: critDmg,
		})

	}

	return r
}

func (a artifacts) substatTotal() map[statTypes]float64 {

	r := make(map[statTypes]float64)
	r[a.Flower.Main.t] += a.Flower.Main.s
	r[a.Feather.Main.t] += a.Feather.Main.s
	r[a.Sands.Main.t] += a.Sands.Main.s
	r[a.Goblet.Main.t] += a.Goblet.Main.s
	r[a.Circlet.Main.t] += a.Circlet.Main.s

	//sub stats
	for _, v := range a.Flower.Sub {
		r[v.t] += v.s
	}
	for _, v := range a.Feather.Sub {
		r[v.t] += v.s
	}
	for _, v := range a.Sands.Sub {
		r[v.t] += v.s
	}
	for _, v := range a.Goblet.Sub {
		r[v.t] += v.s
	}
	for _, v := range a.Circlet.Sub {
		r[v.t] += v.s
	}

	return r
}

func test(p profile) {
	//make our current char and see how accurate we get
	var a artifacts

	a.Flower.Main = stat{t: sHP, s: 4780}
	a.Flower.Sub = []stat{
		stat{
			s: 47,
			t: sATK,
		},
		stat{
			s: 0.179,
			t: sCD,
		},
	}
	a.Feather.Main.t = sATK
	a.Feather.Main.s = 311
	a.Feather.Sub = []stat{
		stat{
			s: 0.053,
			t: sATKP,
		},
		stat{
			s: 0.089,
			t: sCC,
		},
	}
	a.Sands.Main.t = sATKP
	a.Sands.Main.s = 0.466
	a.Sands.Sub = []stat{
		stat{
			s: 31,
			t: sATK,
		},
		stat{
			s: 0.225,
			t: sCD,
		},
	}
	a.Goblet.Main.t = sEleP
	a.Goblet.Main.s = 0.466
	a.Goblet.Sub = []stat{
		stat{
			s: 0.047,
			t: sATKP,
		},
		stat{
			s: 0.097,
			t: sCC,
		},
	}
	a.Circlet.Main.t = sCD
	a.Circlet.Main.s = 0.622
	a.Circlet.Sub = []stat{
		stat{
			s: 0.14,
			t: sATKP,
		},
		stat{
			s: 0.097,
			t: sCC,
		},
	}

	r := calc(a, p)

	for i, v := range r {
		fmt.Printf("Talent %v normal %.0f, crit %.0f, avg %.0f\n", i, v.n, v.c, v.a)
	}
}

func loadMainStatScaling(path string) (map[statTypes][]float64, error) {
	//load substat weights
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	result := make(map[statTypes][]float64)

	if len(lines) < 13 {
		return nil, fmt.Errorf("unexpectedly short main stat scaling file")
	}

	for i := 1; i < len(lines); i++ {
		line := lines[i]

		if len(line) < 22 {
			return nil, fmt.Errorf("unexpectedly short line %v", i)
		}

		for j := 1; j < len(line); j++ {
			val, err := strconv.ParseFloat(line[j], 64)
			if err != nil {
				fmt.Printf("main stat scale - err parsing float at line: %v, pos: %v, line = %v\n", i, j, line[j])
				return nil, err
			}
			result[statTypes(line[0])] = append(result[statTypes(line[0])], val)
		}
	}
	return result, nil
}

func loadSubstatTier(path string) ([]map[statTypes]float64, error) {
	//load substat weights
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	var result []map[statTypes]float64
	//initialize the maps
	for i := 0; i < 4; i++ {
		result = append(result, make(map[statTypes]float64))
	}

	if len(lines) < 2 {
		return nil, fmt.Errorf("unexpectedly short substat tier file")
	}

	for i := 1; i < len(lines); i++ {
		line := lines[i]

		if len(line) < 5 {
			return nil, fmt.Errorf("unexpectedly short line %v", i)
		}

		for j := 1; j < len(line); j++ {
			val, err := strconv.ParseFloat(line[j], 64)
			if err != nil {
				fmt.Printf("substat tier - err parsing float at line: %v, pos: %v, line = %v\n", i, j, line[j])
				return nil, err
			}

			result[j-1][statTypes(line[0])] = val
		}

	}
	return result, nil
}
