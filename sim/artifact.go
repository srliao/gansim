package sim

import "log"

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
func (a *Artifact) Validate() bool {

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
		max := cfg.SubstatTier[3][v.t] * float64(up)
		if v.s > max {
			log.Panicf("invalid feather detected, substat %v exceed 5x max tier. %v\n", v.t, r)
		}
	}

	return true

}
