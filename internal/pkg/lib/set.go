package lib

//Set of 5 artifacts
type Set struct {
	Artifacts map[Slot]Artifact
}

//TotalStats calculate total stats of a set of artifacts
func (s *Set) TotalStats() map[StatType]float64 {
	r := make(map[StatType]float64)

	for _, a := range s.Artifacts {
		r[a.MainStat.Type] += a.MainStat.Value

		for _, v := range a.Substat {
			r[v.Type] += v.Value
		}

	}

	return r
}
