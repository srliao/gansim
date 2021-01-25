package lib

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

func prettySet(s map[Slot]Artifact) string {
	var sb strings.Builder

	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)

	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%v [%v]; ", k, s[Slot(k)].pretty()))
	}

	return sb.String()
}

func prettySetM(s map[Slot]Artifact) string {
	var sb strings.Builder

	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)

	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("\t%v [%v]\n", k, s[Slot(k)].pretty()))
	}

	artifactStats := make(map[StatType]float64)

	for _, a := range s {
		artifactStats[a.MainStat.Type] += a.MainStat.Value

		for _, v := range a.Substat {
			artifactStats[v.Type] += v.Value
		}

	}

	skeys := make([]string, 0, len(artifactStats))
	for k := range artifactStats {
		skeys = append(skeys, string(k))
	}
	sort.Strings(skeys)
	sb.WriteString("\t")
	for _, k := range skeys {
		sb.WriteString(fmt.Sprintf("%v: %.4f; ", k, artifactStats[StatType(k)]))
	}
	sb.WriteString("\n")

	return sb.String()
}

//Artifact represents one artfact
type Artifact struct {
	Level    int64  `yaml:"Level"`
	Type     Slot   `yaml:"Type"`
	MainStat Stat   `yaml:"MainStat"`
	Substat  []Stat `yaml:"Substat"`
}

//TotalStats calculate total stats of one artifact
func (a *Artifact) TotalStats() map[StatType]float64 {
	r := make(map[StatType]float64)

	r[a.MainStat.Type] += a.MainStat.Value

	for _, v := range a.Substat {
		r[v.Type] += v.Value
	}

	return r
}

func (a Artifact) pretty() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("m:%v ", a.MainStat.Type))
	for _, v := range a.Substat {
		sb.WriteString(fmt.Sprintf(" %v:%.4f", v.Type, v.Value))
	}

	return sb.String()
}

//Validate checks if this artifact is valid
func (a *Artifact) Validate(subTier map[int64]map[StatType]float64) bool {

	//no duplicated stats
	dup := make(map[StatType]bool)

	dup[a.MainStat.Type] = true
	count := 0

	for _, v := range a.Substat {
		dup[v.Type] = true
		count++
	}

	if len(dup) != count {
		return false //there's a dup
	}

	//substat amount cannot exceed max tier for the lvl
	up := a.Level / 4

	for _, v := range a.Substat {
		max := subTier[3][v.Type] * float64(up)
		if v.Value > max {
			log.Panicf("invalid feather detected, substat %v exceed %vx max tier. %v\n", v.Type, up, a)
		}
	}

	return true

}
