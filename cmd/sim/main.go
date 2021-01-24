package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/srliao/gansim/internal/pkg/lib"
	"gopkg.in/yaml.v2"
)

type config struct {
	Profiles     []string `yaml:"Profiles"`
	GraphOutput  string   `yaml:"GraphOutput"`
	NumSim       int64    `yaml:"NumSim"`
	BinSize      int64    `yaml:"BinSize"`
	WriteHist    bool     `yaml:"WriteHist"`
	DamageType   string   `yaml:"DamageType"`
	NumWorker    int64    `yaml:"NumWorker"`
	MainStatFile string   `yaml:"MainStatFile"`
	SubTierFile  string   `yaml:"SubstatTierFile"`
	MainProbFile string   `yaml:"MainStatProbFile"`
	SubProbFile  string   `yaml:"SubProbFile"`
}

type artifactConfig struct {
	MainStats map[lib.Slot]lib.StatType `yaml:"MainStats"`
	Level     int64                     `yaml:"Level"`
}

type simProfile struct {
	Profile          lib.Profile               `yaml:"Profile"`
	ArtifactConfig   artifactConfig            `yaml:"ArtifactConfig"`
	CurrentArtifacts map[lib.Slot]lib.Artifact `yaml:"CurrentArtifacts"`
}

func main() {

	var src []byte
	var cfg config
	var err error

	src, err = ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(src, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = sim(cfg)

	if err != nil {
		log.Fatal(err)
	}

}

func sim(cfg config) error {
	//load each profile

	var src []byte
	var err error
	var p simProfile

	ms, err := loadMainStat(cfg.MainStatFile)
	if err != nil {
		return err
	}
	fmt.Println("main stat scaling loaded ok")
	// fmt.Println(ms)
	mp, err := loadMainProb(cfg.MainProbFile)
	if err != nil {
		return err
	}
	fmt.Println("main stat prob loaded ok")
	// fmt.Println(mp)
	st, err := loadSubTier(cfg.SubTierFile)
	if err != nil {
		return err
	}
	fmt.Println("sub tiers loaded ok")
	// fmt.Println(st)
	sp, err := loadSubProb(cfg.SubProbFile)
	if err != nil {
		return err
	}
	fmt.Println("substat prob loaded ok")
	// fmt.Println(sp)

	histData := make([][]float64, len(cfg.Profiles))
	labels := make([]string, len(cfg.Profiles))

	for i, pp := range cfg.Profiles {
		src, err = ioutil.ReadFile(pp)
		if err != nil {
			fmt.Println("error reading file")
			return err
		}
		err = yaml.Unmarshal(src, &p)
		if err != nil {
			return err
		}

		if _, ok := p.ArtifactConfig.MainStats[lib.Sands]; !ok {
			return fmt.Errorf("invalid profile: no stats specified for sands")
		}
		if _, ok := p.ArtifactConfig.MainStats[lib.Goblet]; !ok {
			return fmt.Errorf("invalid profile: no stats specified for goblet")
		}
		if _, ok := p.ArtifactConfig.MainStats[lib.Circlet]; !ok {
			return fmt.Errorf("invalid profile: no stats specified for circlet")
		}

		labels[i] = fmt.Sprintf("%v [lvl %v]", p.Profile.Label, p.ArtifactConfig.Level)

		//calculate the damage distribution
		resp := make(chan lib.DmgResult, cfg.NumSim)

		var hist []float64

		fmt.Printf("starting simulation for profile: %v, n = %v\n", pp, cfg.NumSim)
		timeStart := time.Now()

		//fire up max number of workers
		req := make(chan bool)
		done := make(chan bool)
		for w := 0; w < int(cfg.NumWorker); w++ {
			g, err := lib.NewGenerator(
				time.Now().UnixNano()+int64(w),
				ms,
				mp,
				st,
				sp,
			)
			if err != nil {
				return err
			}
			go workerD(g, p, resp, req, done)
		}

		//use a go routine to send out a job whenever a worker is done
		go func() {
			var wip int64
			for wip < cfg.NumSim {
				//try sending a job to req chan while wip < cfg.NumSim
				req <- true
				wip++
			}
		}()

		count := cfg.NumSim
		var progress float64
		fmt.Print("\tProgress: 0%")
		for count > 0 {
			//process results received
			r := <-resp
			count--

			val := r.Avg
			switch cfg.DamageType {
			case "normal":
				val = r.Normal
			case "crit":
				val = r.Crit
			}

			//push result to r
			hist = append(hist, val)

			if (1 - float64(count)/float64(cfg.NumSim)) > (progress + 0.1) {
				progress = (1 - float64(count)/float64(cfg.NumSim))
				fmt.Printf("...%.0f%%", 100*progress)
			}
		}
		fmt.Print("...100%%\n")

		close(done)

		elapsed := time.Since(timeStart)
		fmt.Printf("Simulation for profile %v took %s\n\n", pp, elapsed)

		histData[i] = hist

		//figure out damage required to hit required percentile

		//sim distribution to reach said dmg
	}

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

	return nil
}

func workerD(g *lib.Generator, p simProfile, resp chan lib.DmgResult, req chan bool, done chan bool) {
	for {
		select {
		case <-req:
		case <-done:
			return
		}
		//generate a set of artifacts
		set := make(map[lib.Slot]lib.Artifact)
		set[lib.Flower] = g.RandWithMain(lib.Flower, lib.HP, p.ArtifactConfig.Level)
		set[lib.Feather] = g.RandWithMain(lib.Feather, lib.ATK, p.ArtifactConfig.Level)
		set[lib.Sands] = g.RandWithMain(lib.Sands, p.ArtifactConfig.MainStats[lib.Sands], p.ArtifactConfig.Level)
		set[lib.Goblet] = g.RandWithMain(lib.Goblet, p.ArtifactConfig.MainStats[lib.Goblet], p.ArtifactConfig.Level)
		set[lib.Circlet] = g.RandWithMain(lib.Circlet, p.ArtifactConfig.MainStats[lib.Circlet], p.ArtifactConfig.Level)

		//calculate dmg
		r := lib.Calc(p.Profile, set, false)

		var out lib.DmgResult
		for _, v := range r {
			out.Normal += v.Normal
			out.Avg += v.Avg
			out.Crit += v.Crit
		}

		resp <- out
	}
}

func loadMainStat(path string) (map[lib.StatType][]float64, error) {
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
	result := make(map[lib.StatType][]float64)

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
			result[lib.StatType(line[0])] = append(result[lib.StatType(line[0])], val)
		}
	}
	return result, nil
}

func loadMainProb(path string) (map[lib.Slot][]lib.StatProb, error) {
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
	result := make(map[lib.Slot][]lib.StatProb)

	//read header

	if len(lines) < 1 {
		return nil, fmt.Errorf("unexpected short file")
	}

	if len(lines[0]) < 6 {
		return nil, fmt.Errorf("unexpected short header line")
	}

	var header []lib.Slot

	for i := 1; i < len(lines[0]); i++ {
		result[lib.Slot(lines[0][i])] = make([]lib.StatProb, 0)
		header = append(header, lib.Slot(lines[0][i]))
	}

	// fmt.Println(header)

	for i := 1; i < len(lines); i++ {
		if len(lines[i]) != 6 {
			return nil, fmt.Errorf("line %v does not have 6 fields, got %v", i, len(lines[i]))
		}
		for j := 1; j < len(lines[i]); j++ {
			prb, err := strconv.ParseFloat(lines[i][j], 64)
			if err != nil {
				return nil, fmt.Errorf("err parsing float @ line %v, value %v: %v", i, lines[i][j], err)
			}
			result[header[j-1]] = append(
				result[header[j-1]],
				lib.StatProb{
					Type: lib.StatType(lines[i][0]),
					Prob: prb,
				},
			)

		}
	}

	return result, nil
}

func loadSubTier(path string) ([]map[lib.StatType]float64, error) {
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
	var result []map[lib.StatType]float64
	//initialize the maps
	for i := 0; i < 4; i++ {
		result = append(result, make(map[lib.StatType]float64))
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

			result[j-1][lib.StatType(line[0])] = val
		}

	}
	return result, nil
}

func loadSubProb(path string) (map[lib.StatType][]lib.StatProb, error) {
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
	result := make(map[lib.StatType][]lib.StatProb)

	//read header

	if len(lines) < 1 {
		return nil, fmt.Errorf("unexpected short file")
	}

	if len(lines[0]) < 11 {
		return nil, fmt.Errorf("unexpected short header line")
	}

	var header []lib.StatType

	for i := 1; i < len(lines[0]); i++ {
		result[lib.StatType(lines[0][i])] = make([]lib.StatProb, 0)
		header = append(header, lib.StatType(lines[0][i]))
	}

	// fmt.Println(header)

	for i := 1; i < len(lines); i++ {
		if len(lines[i]) != 11 {
			return nil, fmt.Errorf("line %v does not have 11 fields, got %v", i, len(lines[i]))
		}
		for j := 1; j < len(lines[i]); j++ {
			prb, err := strconv.ParseFloat(lines[i][j], 64)
			if err != nil {
				return nil, fmt.Errorf("err parsing float @ line %v, value %v: %v", i, lines[i][j], err)
			}
			result[header[j-1]] = append(
				result[header[j-1]],
				lib.StatProb{
					Type: lib.StatType(lines[i][0]),
					Prob: prb,
				},
			)

		}
	}

	return result, nil
}

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
