package main

import (
	"encoding/csv"
	"fmt"
	"io"
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
	Profiles    []string `yaml:"Profiles"`
	GraphOutput string   `yaml:"GraphOutput"`
	NumSim      int64    `yaml:"NumSim"`
	BinSize     int64    `yaml:"BinSize"`
	WriteCSV    bool     `yaml:"WriteCSV"`
	DamageType  string   `yaml:"DamageType"`
}

type profile struct {
	NumSim  int64  `yaml:"NumSim"`
	Output  string `yaml:"Output"`
	Label   string `yaml:"Label"`
	SaveCSV bool   `yaml:"SaveCSV"`
	//character info
	CharLevel     float64 `yaml:"CharacterLevel"`
	CharBaseAtk   float64 `yaml:"CharacterBaseAtk"`
	WeaponBaseAtk float64 `yaml:"WeaponBaseAtk"`
	EnemyLevel    float64 `yaml:"EnemyLevel"`
	//artifact info
	Sands          statTypes `yaml:"Sands"`
	Goblet         statTypes `yaml:"Goblet"`
	Circlet        statTypes `yaml:"Circlet"`
	SandsStat      float64   `yaml:"SandsStat"`
	GobletStat     float64   `yaml:"GobletStat"`
	CircletStat    float64   `yaml:"CircletStat"`
	SubstatFile    string    `yaml:"SubstatFile"`
	SubstatWeights map[string][]statPrb
	//talents list
	Talents []float64 `yaml:"Talents"`
	//stat mods
	AtkMod      [][]float64 `yaml:"AtkMod"`
	EleMod      [][]float64 `yaml:"EleMod"`
	CCMod       [][]float64 `yaml:"CCMod"`
	CDMod       [][]float64 `yaml:"CDMod"`
	DmgMod      [][]float64 `yaml:"DmgMod"`
	ResistMod   [][]float64 `yaml:"ResistMod"`
	DefShredMod [][]float64 `yaml:"DefShredMod"`
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

var subtier []map[statTypes]float64

func main() {
	// runtime.GOMAXPROCS(12)
	//read config

	var err error

	f, err := os.Open("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// log.Println("reading config file")
	var cfg config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(cfg)

	//initialize stat maps
	t1 := make(map[statTypes]float64)
	t1[sHP] = 209
	t1[sDEF] = 16
	t1[sATK] = 14
	t1[sHPP] = 0.041
	t1[sDEFP] = 0.051
	t1[sATKP] = 0.041
	t1[sEM] = 16
	t1[sER] = 0.045
	t1[sCC] = 0.027
	t1[sCD] = 0.054
	t2 := make(map[statTypes]float64)
	t2[sHP] = 239
	t2[sDEF] = 19
	t2[sATK] = 16
	t2[sHPP] = 0.047
	t2[sDEFP] = 0.058
	t2[sATKP] = 0.047
	t2[sEM] = 19
	t2[sER] = 0.052
	t2[sCC] = 0.031
	t2[sCD] = 0.062
	t3 := make(map[statTypes]float64)
	t3[sHP] = 269
	t3[sDEF] = 21
	t3[sATK] = 18
	t3[sHPP] = 0.053
	t3[sDEFP] = 0.066
	t3[sATKP] = 0.053
	t3[sEM] = 21
	t3[sER] = 0.058
	t3[sCC] = 0.035
	t3[sCD] = 0.07
	t4 := make(map[statTypes]float64)
	t4[sHP] = 299
	t4[sDEF] = 23
	t4[sATK] = 19
	t4[sHPP] = 0.058
	t4[sDEFP] = 0.073
	t4[sATKP] = 0.058
	t4[sEM] = 23
	t4[sER] = 0.065
	t4[sCC] = 0.039
	t4[sCD] = 0.078

	subtier = append(subtier, t1)
	subtier = append(subtier, t2)
	subtier = append(subtier, t3)
	subtier = append(subtier, t4)

	rand.Seed(time.Now().UTC().UnixNano())

	//loop through profiles and run sim for each

	histData := make([][]float64, len(cfg.Profiles))
	labels := make([]string, len(cfg.Profiles))

	for i, ppath := range cfg.Profiles {
		fp, err := os.Open(ppath)
		if err != nil {
			log.Fatal(err)
		}
		defer fp.Close()

		// log.Printf("reading profile: %v\n", ppath)
		var prf profile
		pdecoder := yaml.NewDecoder(fp)
		err = pdecoder.Decode(&prf)
		if err != nil {
			log.Fatal(err)
		}

		//profile sanity check
		if len(prf.AtkMod) < len(prf.Talents) {
			log.Panicln("invalid # of AtkMod")
		}
		if len(prf.EleMod) < len(prf.Talents) {
			log.Panicln("invalid # of EleMod")
		}
		if len(prf.CCMod) < len(prf.Talents) {
			log.Panicln("invalid # of CCMod")
		}
		if len(prf.CDMod) < len(prf.Talents) {
			log.Panicln("invalid # of CDMod")
		}
		if len(prf.DmgMod) < len(prf.Talents) {
			log.Panicln("invalid # of DmgMod")
		}
		if len(prf.ResistMod) < len(prf.Talents) {
			log.Panicln("invalid # of ResistMod")
		}
		if len(prf.DefShredMod) < len(prf.Talents) {
			log.Panicln("invalid # of DefShredMod")
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

		labels[i] = prf.Label

		fmt.Printf("starting simulation for profile: %v, n = %v\n", ppath, prf.NumSim)
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
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
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

	wc := 0
	wcmax := 12
	resp := make(chan result)

	// fmt.Println("n = ", count)
	//keep sending out workers while still simulations left to do
	fmt.Print("\tProgress: 0%")
	for count > 0 {
		//send out worker if wc < wcmax
		if wc < wcmax {
			go worker(p, resp)
			wc++
		} else { //otherwise wait for result
			r := <-resp
			wc--
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

	}
	fmt.Print("...100%%\n")

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

func worker(p profile, resp chan result) {
	art := genArtifacts(p)
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

//generate a set of artifacts given the configs
func genArtifacts(p profile) artifacts {
	var r artifacts
	//flower is always hp = 4780
	r.Flower = randArtifact(stat{t: sHP, s: 4780}, p.SubstatWeights["flower"])
	//feather is always flat atk 311
	r.Feather = randArtifact(stat{t: sATK, s: 311}, p.SubstatWeights["feather"])
	//sands is always % atk 46.6%
	r.Sands = randArtifact(stat{t: p.Sands, s: p.SandsStat}, p.SubstatWeights["sands"])
	//goblet is always % ele 46.6%
	r.Goblet = randArtifact(stat{t: p.Goblet, s: p.GobletStat}, p.SubstatWeights["goblet"])
	//circlet is always crit dmg 62.20%
	r.Circlet = randArtifact(stat{t: p.Circlet, s: p.CircletStat}, p.SubstatWeights["circlet"])
	//do some sanity checks on sub stat
	for _, v := range r.Flower.Sub {
		max := subtier[3][v.t] * 5
		if v.s > max {
			log.Panicf("invalid flower detected, substat %v exceed 5x max tier. %v\n", v.t, r)
		}
		if v.t == r.Flower.Main.t {
			log.Panicf("invalid flower detected, substat %v is same as main stat. %v\n", v.t, r)
		}
	}
	for _, v := range r.Feather.Sub {
		max := subtier[3][v.t] * 5
		if v.s > max {
			log.Panicf("invalid feather detected, substat %v exceed 5x max tier. %v\n", v.t, r)
		}
		if v.t == r.Feather.Main.t {
			log.Panicf("invalid feather detected, substat %v is same as main stat. %v\n", v.t, r)
		}
	}
	for _, v := range r.Sands.Sub {
		max := subtier[3][v.t] * 5
		if v.s > max {
			log.Panicf("invalid sands detected, substat %v exceed 5x max tier. %v\n", v.t, r)
		}
		if v.t == r.Sands.Main.t {
			log.Panicf("invalid sands detected, substat %v is same as main stat. %v\n", v.t, r)
		}
	}
	for _, v := range r.Goblet.Sub {
		max := subtier[3][v.t] * 5
		if v.s > max {
			log.Panicf("invalid goblet detected, substat %v exceed 5x max tier. %v\n", v.t, r)
		}
		if v.t == r.Goblet.Main.t {
			log.Panicf("invalid goblet detected, substat %v is same as main stat. %v\n", v.t, r)
		}
	}
	for _, v := range r.Circlet.Sub {
		max := subtier[3][v.t] * 5
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
func randArtifact(main stat, prb []statPrb) artifact {
	var r artifact
	r.Main.s = main.s
	r.Main.t = main.t
	//how many substats
	var prb4stat, sum, nextProbSum float64
	var next []statPrb
	var found bool
	prb4stat = 207 / 932
	p := rand.Float64()
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

	for i := 0; i < 4; i++ {
		var current []statPrb
		var check float64
		for _, v := range next {
			current = append(current, statPrb{t: v.t, p: v.p / nextProbSum})
			check += v.p / nextProbSum
		}
		if showDebug {
			log.Println("current probabilities: ", current)
			log.Println("sub stat count: ", len(current))
			log.Println("current prob total: ", check)
		}
		p = rand.Float64()
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
				tier := rand.Intn(4)
				val := subtier[tier][v.t]
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

	if len(r.Sub) != 4 {
		log.Panicln("expected to have 4 substat lines here, got ", len(r.Sub))
	}

	//if line == 4, then upgrade once, otherwise skip since first roll will be the 4th line
	if lines == 4 {
		//upgrade
		//ASSUMPTION EQUAL CHANCE OF UPGRADING EACH STAT PROB NOT TRUE???
		i := rand.Intn(4)
		tier := rand.Intn(4)
		r.Sub[i].s += subtier[tier][r.Sub[i].t]

	}

	//+8
	//ASSUMPTION EQUAL CHANCE OF UPGRADING EACH STAT PROB NOT TRUE???
	i := rand.Intn(4)
	tier := rand.Intn(4)
	r.Sub[i].s += subtier[tier][r.Sub[i].t]

	//+12
	//ASSUMPTION EQUAL CHANCE OF UPGRADING EACH STAT PROB NOT TRUE???
	i = rand.Intn(4)
	tier = rand.Intn(4)
	r.Sub[i].s += subtier[tier][r.Sub[i].t]

	//+16
	//ASSUMPTION EQUAL CHANCE OF UPGRADING EACH STAT PROB NOT TRUE???
	i = rand.Intn(4)
	tier = rand.Intn(4)
	r.Sub[i].s += subtier[tier][r.Sub[i].t]

	//+20
	//ASSUMPTION EQUAL CHANCE OF UPGRADING EACH STAT PROB NOT TRUE???
	i = rand.Intn(4)
	tier = rand.Intn(4)
	r.Sub[i].s += subtier[tier][r.Sub[i].t]

	return r
}

func calc(a artifacts, p profile) []result {

	//artifact substats
	artSubStat := a.substatTotal()

	var r []result

	//loop through each talent
	for i, t := range p.Talents {
		//calculate total atk
		var totalAtk, atkp, cc, cd, atk, elep, dmgBonus, defAdj, resist float64
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
		for _, v := range p.DefShredMod[i] {
			defShred += v
		}
		//calculate def adjustment
		defAdj = (100 + p.CharLevel) / ((100 + p.CharLevel) + (100+p.EnemyLevel)*(1-defShred))

		//add up atk % mods
		for _, v := range p.AtkMod[i] {
			atkp += v
		}

		totalAtk = totalAtk*(1+atkp) + atk

		//add up dmg mods
		for _, v := range p.EleMod[i] {
			dmgBonus += v
		}
		for _, v := range p.DmgMod[i] {
			dmgBonus += v
		}

		dmgBonus += elep //add in ele bonus from artifacts

		//add up crit mods
		for _, v := range p.CCMod[i] {
			cc += v
		}
		for _, v := range p.CDMod[i] {
			cd += v
		}

		//calculate enemy resistance
		for _, v := range p.ResistMod[i] {
			//if v will bring resist to negative, the half the portion that brings it to negative
			if resist+v < 0 {
				//if resist is already negative
				if resist < 0 {
					//if v is positive?? just don't list it this way
					if v >= 0 {
						resist += v
					} else {
						//half the effect of v
						resist += v / 2
					}
				} else {
					temp := v + resist
					resist = 0
					resist += temp / 2
				}
			} else {
				resist += v
			}
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

		normalDmg := totalAtk * (1 + dmgBonus) * t * defAdj * (1 - resist)
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
