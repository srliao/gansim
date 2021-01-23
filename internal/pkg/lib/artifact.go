package lib

import "log"

//Artifact represents one artfact
type Artifact struct {
	Level    int64
	Type     Slot
	MainStat Stat
	Substat  []Stat
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
