package sim

//Slot identifies the artifact slot
type Slot string

//Types of artifact slots
const (
	Flower  Slot = "Flower"
	Feather Slot = "Feather"
	Sands   Slot = "Sands"
	Goblet  Slot = "Goblet"
	Circlet Slot = "Circlet"
)

type Stat struct {
	Type  StatType
	Value float64
}

type StatType string

//stat types
const (
	DEFP StatType = "DEF%"
	DEF  StatType = "DEF"
	HP   StatType = "HP"
	HPP  StatType = "HP%"
	ATK  StatType = "ATK"
	ATKP StatType = "ATK%"
	ER   StatType = "ER"
	EM   StatType = "EM"
	CC   StatType = "CR"
	CD   StatType = "CD"
	Heal StatType = "Heal"
	EleP StatType = "Ele%"
	PhyP StatType = "Phys%"
)
