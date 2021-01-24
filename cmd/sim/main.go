package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

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

	ms, err := loadMainStat(cfg.MainStatFile)
	if err != nil {
		log.Fatal(err)
	}
	st, err := loadSubTier(cfg.SubTierFile)
	if err != nil {
		log.Fatal(err)
	}

	g, err := lib.NewGenerator(
		time.Now().UnixNano(),
		ms,
		nil,
		st,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(g)

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
	var result map[lib.Slot][]lib.StatProb

	//read header

	if len(lines) < 1 {
		return nil, fmt.Errorf("unexpected short file")
	}

	if len(lines[0]) < 6 {
		return nil, fmt.Errorf("unexpected short header line")
	}

	var header []lib.Slot

	for i := 1; i < len(lines[0]); i++ {
		result[lib.Slot(lines[0][i])] = make([]lib.StatProb, 5)
		header = append(header, lib.Slot(lines[0][i]))
	}

	for i := 1; i < len(lines); i++ {
		if len(lines[i]) != 6 {
			return nil, fmt.Errorf("line %v does not have 6 fields, got %v", i, len(lines[i]))
		}
		for j := 1; j < len(lines[i]); j++ {
			result[header[j-1])][j-1] = lines[i][j]

		}
	}

	return nil, nil
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
