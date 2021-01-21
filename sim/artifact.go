package sim

//set of 5 artifacts
type Set struct {
}

type Slot string

//Types of artifact slots
const (
	Flower  Slot = "Flower"
	Feather Slot = "Feather"
	Sands   Slot = "Sands"
	Goblet  Slot = "Goblet"
	Circlet Slot = "Circlet"
)

type Artifact struct {
	Stat  StatType
	Value float64
}

type StatType string

const (
	//stat types
	sDEFP StatType = "DEF%"
	sDEF  StatType = "DEF"
	sHP   StatType = "HP"
	sHPP  StatType = "HP%"
	sATK  StatType = "ATK"
	sATKP StatType = "ATK%"
	sER   StatType = "ER"
	sEM   StatType = "EM"
	sCC   StatType = "CR"
	sCD   StatType = "CD"
	sHeal StatType = "Heal"
	sEleP StatType = "Ele%"
	sPhyP StatType = "Phys%"
)
