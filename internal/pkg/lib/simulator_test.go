package lib

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func testGenerator(t *testing.T) {

	ms, err := loadMainStat("./test/mainstat.csv")
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Println(ms)
	mp, err := loadMainProb("./test/mainprob.csv")
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(mp)
	st, err := loadSubTier("./test/substat.csv")
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(st)
	sp, err := loadSubProb("./test/subprob.csv")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(sp[Flower])

	//make sure substat prob converges
	s, err := NewSimulator(
		ms,
		mp,
		st,
		sp,
		false,
		false,
	)

	seed := time.Now().UnixNano()
	rand := rand.New(rand.NewSource(seed))

	fmt.Println(sp[Flower])

	//map out the substat probs
	expected := make(map[Slot]map[StatType][]StatProb)

	for slot, x := range sp {
		expected[slot] = make(map[StatType][]StatProb)
		for mainST, y := range x {

			total := 0.0

			expected[slot][mainST] = y

			for _, v := range y {
				total += v.Prob
			}

			var check float64

			for i, v := range expected[slot][mainST] {
				expected[slot][mainST][i].Prob = v.Prob / total
				check += v.Prob / total
			}

			fmt.Printf("Slot %v, %v, check = %v\n", slot, mainST, check)
		}
	}

	testCount := 1000000

	for rs := 0; rs < 5; rs++ {
		t := Circlet
		switch rs {
		case 0:
			t = Flower
		case 1:
			t = Feather
		case 2:
			t = Sands
		case 3:
			t = Goblet

		}

		//loop through each slot
		for _, mainStat := range mp[t] {

			if mainStat.Prob == 0 {
				continue
			}

			fmt.Printf("%v %v: \n", t, mainStat.Type)

			prob := make(map[StatType]int64)
			var total int64
			var noatk int64

			//generate sub stats given mainStat 100000 and check if probability converges
			for i := 0; i < testCount; i++ {
				next := s.RandArtifact(t, mainStat.Type, 20, rand)

				hasAttack := false
				for _, sub := range next.Substat {
					prob[sub.Type]++
					total++

					if sub.Type == ATK || sub.Type == ATKP || sub.Type == CR || sub.Type == CD {
						hasAttack = true
					}
				}

				if !hasAttack {
					noatk++
				}
			}

			for _, v := range expected[t][mainStat.Type] {
				fmt.Printf("\tat %v expected %.5f, got %.5f, diff %.5f\n", v.Type, v.Prob, float64(prob[v.Type])/float64(total), v.Prob-float64(prob[v.Type])/float64(total))

			}

			fmt.Printf("\tno attack %v\n", noatk)

			// fmt.Println(prob)
		}
	}

}

func loadMainStat(path string) (map[StatType][]float64, error) {
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
	result := make(map[StatType][]float64)

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
			result[StatType(line[0])] = append(result[StatType(line[0])], val)
		}
	}
	return result, nil
}

func loadMainProb(path string) (map[Slot][]StatProb, error) {
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
	result := make(map[Slot][]StatProb)

	//read header

	if len(lines) < 1 {
		return nil, fmt.Errorf("unexpected short file")
	}

	if len(lines[0]) < 6 {
		return nil, fmt.Errorf("unexpected short header line")
	}

	var header []Slot

	for i := 1; i < len(lines[0]); i++ {
		result[Slot(lines[0][i])] = make([]StatProb, 0)
		header = append(header, Slot(lines[0][i]))
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
				StatProb{
					Type: StatType(lines[i][0]),
					Prob: prb,
				},
			)

		}
	}

	return result, nil
}

func loadSubTier(path string) ([]map[StatType]float64, error) {
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
	var result []map[StatType]float64
	//initialize the maps
	for i := 0; i < 4; i++ {
		result = append(result, make(map[StatType]float64))
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

			result[j-1][StatType(line[0])] = val
		}

	}
	return result, nil
}

func loadSubProb(path string) (map[Slot]map[StatType][]StatProb, error) {
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
	result := make(map[Slot]map[StatType][]StatProb)

	//read header

	if len(lines) < 1 {
		return nil, fmt.Errorf("unexpected short file")
	}

	if len(lines[0]) < 12 {
		return nil, fmt.Errorf("unexpected short header line")
	}

	var header []StatType

	for i := 2; i < len(lines[0]); i++ {
		header = append(header, StatType(lines[0][i]))
	}

	// fmt.Println(header)

	for i := 1; i < len(lines); i++ {
		l := lines[i]

		if len(l) != 12 {
			return nil, fmt.Errorf("line %v does not have 12 fields, got %v", i, len(lines[i]))
		}

		slot := Slot(l[0])
		main := StatType(l[1])

		if _, ok := result[slot]; !ok {
			result[slot] = make(map[StatType][]StatProb)
		}
		if _, ok := result[slot][main]; !ok {
			result[slot][main] = make([]StatProb, 0)
		}

		for j := 2; j < len(lines[i]); j++ {
			prb, err := strconv.ParseFloat(l[j], 64)
			if err != nil {
				return nil, fmt.Errorf("err parsing float @ line %v, value %v: %v", i, lines[i][j], err)
			}
			result[slot][main] = append(
				result[slot][main],
				StatProb{
					Type: header[j-2],
					Prob: prb,
				},
			)

		}
	}

	return result, nil
}
