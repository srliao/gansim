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
	DmgBinSize   int64    `yaml:"DmgBinSize"`
	FarmBinSize  int64    `yaml:"FarmBinSize"`
	WriteHist    bool     `yaml:"WriteHist"`
	NumWorker    int64    `yaml:"NumWorker"`
	Percentile   float64  `yaml:"Percentile"`
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

	log.Println(cfg)

	ms, err := loadMainStat(cfg.MainStatFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("main stat scaling loaded ok")
	// fmt.Println(ms)
	mp, err := loadMainProb(cfg.MainProbFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("main stat prob loaded ok")
	// fmt.Println(mp)
	st, err := loadSubTier(cfg.SubTierFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("sub tiers loaded ok")
	// fmt.Println(st)
	sp, err := loadSubProb(cfg.SubProbFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("substat prob loaded ok")
	// fmt.Println(sp)

	labels := make([]string, len(cfg.Profiles))
	page := components.NewPage()
	page.PageTitle = "simulation results"

	histData := make([][]float64, len(cfg.Profiles))
	histStart := make([]int64, len(cfg.Profiles))
	min := make([]float64, len(cfg.Profiles))
	max := make([]float64, len(cfg.Profiles))
	mean := make([]float64, len(cfg.Profiles))
	sd := make([]float64, len(cfg.Profiles))

	//farm charts
	var fcharts []*charts.Line

	//absolute max/min
	var binMin int64
	var binMax int64
	binMin = math.MaxInt64
	binMin = -1

	//load each profile
	var p lib.Profile

	for i, prf := range cfg.Profiles {
		src, err = ioutil.ReadFile(prf)
		if err != nil {
			fmt.Println("error reading file")
			log.Fatal(err)
		}
		err = yaml.Unmarshal(src, &p)
		if err != nil {
			log.Fatal(err)
		}

		if _, ok := p.Artifacts.TargetMainStat[lib.Sands]; !ok {
			log.Fatal("invalid profile: no stats specified for sands")
		}
		if _, ok := p.Artifacts.TargetMainStat[lib.Goblet]; !ok {
			log.Fatal("invalid profile: no stats specified for goblet")
		}
		if _, ok := p.Artifacts.TargetMainStat[lib.Circlet]; !ok {
			log.Fatal("invalid profile: no stats specified for circlet")
		}

		labels[i] = fmt.Sprintf("%v", p.Label)

		s, err := lib.NewSimulator(
			ms,
			mp,
			st,
			sp,
		)
		if err != nil {
			log.Fatal(err)
		}

		ds, dhist, dmin, dmax, dmean, dsd := s.SimDmgDist(cfg.NumSim, cfg.DmgBinSize, cfg.NumWorker, p)

		histData[i] = dhist
		histStart[i] = ds
		min[i] = dmin
		max[i] = dmax
		mean[i] = dmean
		sd[i] = dsd

		if ds < binMin {
			binMin = ds
		}
		m := ds + int64(len(dhist)+1)*cfg.DmgBinSize
		if m > binMax {
			binMax = m
		}

		//figure out damage required to hit required percentile
		var total, cumul, d float64

		for _, v := range dhist {
			total += v
		}

		for i, v := range dhist {
			cumul += v
			if cumul/total >= cfg.Percentile {
				d = float64(i)
			}
		}

		d = float64(ds) + d*float64(cfg.DmgBinSize)

		//sim distribution to reach said dmg
		fstart, fhist, fmin, fmax, fmean, fsd := s.SimArtifactFarm(cfg.NumSim, cfg.FarmBinSize, cfg.NumWorker, d, p)

		var fx []int64
		var fitems []opts.LineData

		for i, v := range fhist {
			fx = append(fx, fstart+cfg.FarmBinSize*int64(i))
			fitems = append(fitems, opts.LineData{Value: v})
		}

		//one chart for every one of these sims
		lineChart := charts.NewLine()
		lineChart.SetGlobalOptions(
			charts.WithTitleOpts(opts.Title{
				Title: fmt.Sprintf("Histogram (n = %v, %.2f percentile)", cfg.NumSim, cfg.Percentile),
			}),
			charts.WithYAxisOpts(opts.YAxis{
				Name: "Freq",
			}),
			charts.WithXAxisOpts(opts.XAxis{
				Name: "# of Artifacts",
			}),
			// charts.WithTooltipOpts(opts.Tooltip{Show: true}),
			charts.WithLegendOpts(opts.Legend{Show: true, Right: "0%", Orient: "vertical", Data: labels}),
		)
		lineChart.AddSeries(labels[i], fitems)
		lineChart.SetXAxis(fx)
		fcharts = append(fcharts, lineChart)

		fmt.Printf("min: %v, max %v, mean: %.2f, sd: %.2f\n", fmin, fmax, fmean, fsd)
	}

	numBin := (binMax - binMin) / cfg.DmgBinSize
	xaxis := make([]float64, numBin)

	for i := range xaxis {
		xaxis[i] = float64(int64(i)*cfg.DmgBinSize + binMin)
	}

	bins := make([][]float64, len(cfg.Profiles))
	items := make([][]opts.LineData, len(cfg.Profiles))

	for i, hist := range histData {
		bins[i] = make([]float64, numBin)
		offset := (histStart[i] - binMin) / cfg.DmgBinSize
		for j, v := range hist {
			bins[i][int(offset)+j] += v
		}
		labels[i] = fmt.Sprintf("%v (min: %.f max: %.f avg: %.f sd: %.f)", labels[i], min[i], max[i], mean[i], sd[i])
	}

	for i, b := range bins {
		for _, v := range b {
			items[i] = append(items[i], opts.LineData{Value: v / float64(cfg.NumSim)})
		}
	}

	lineChart := charts.NewLine()
	lineChart.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: fmt.Sprintf("Probability Density Function (n = %v)", cfg.NumSim),
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Probability",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Dmg",
		}),
		// charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithLegendOpts(opts.Legend{Show: true, Right: "0%", Orient: "vertical", Data: labels}),
	)
	lineChart.SetXAxis(xaxis)

	for i, series := range items {
		lineChart.AddSeries(labels[i], series)
	}

	page.AddCharts(
		lineChart,
	)

	for _, v := range fcharts {
		page.AddCharts(
			v,
		)
	}

	graph, err := os.Create(cfg.GraphOutput)
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(graph))

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

func loadSubProb(path string) (map[lib.Slot]map[lib.StatType][]lib.StatProb, error) {
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
	result := make(map[lib.Slot]map[lib.StatType][]lib.StatProb)

	//read header

	if len(lines) < 1 {
		return nil, fmt.Errorf("unexpected short file")
	}

	if len(lines[0]) < 12 {
		return nil, fmt.Errorf("unexpected short header line")
	}

	var header []lib.StatType

	for i := 2; i < len(lines[0]); i++ {
		header = append(header, lib.StatType(lines[0][i]))
	}

	// fmt.Println(header)

	for i := 1; i < len(lines); i++ {
		l := lines[i]

		if len(l) != 12 {
			return nil, fmt.Errorf("line %v does not have 12 fields, got %v", i, len(lines[i]))
		}

		slot := lib.Slot(l[0])
		main := lib.StatType(l[1])

		if _, ok := result[slot]; !ok {
			result[slot] = make(map[lib.StatType][]lib.StatProb)
		}
		if _, ok := result[slot][main]; !ok {
			result[slot][main] = make([]lib.StatProb, 0)
		}

		for j := 2; j < len(lines[i]); j++ {
			prb, err := strconv.ParseFloat(l[j], 64)
			if err != nil {
				return nil, fmt.Errorf("err parsing float @ line %v, value %v: %v", i, lines[i][j], err)
			}
			result[slot][main] = append(
				result[slot][main],
				lib.StatProb{
					Type: header[j-2],
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
