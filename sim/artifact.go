package sim

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

	//no duplicated substat

	//substat cannot be same as main stat

	//substat amount cannot exceed max tier for the lvl

}
