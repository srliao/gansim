package main

//Profile describe a damage profile to calculate
type profile struct {
	ID                int64             `json:"id"`
	Label             string            `json:"label"`
	Character         character         `json:"character"`
	Weapon            weapon            `json:"weapon"`
	Enemy             enemy             `json:"enemy"`
	ArtifactLvl       int64             `json:"artifact_levels"`
	ArtifactMainStats map[slot]statType `json:"artifact_main_stats"`
	Abilities         []ability         `json:"abilities"`
}

type enemy struct {
	Level   int64   `json:"level"`
	PhysRes float64 `json:"phy_resist"`
	EleRes  float64 `json:"ele_resist"`
}

type character struct {
	Level   int64      `json:"level"`
	BaseAtk float64    `json:"base_atk"`
	Mods    []modifier `json:"mods"`
}

type weapon struct {
	BaseAtk float64    `json:"base_atk"`
	Mods    []modifier `json:"mods"`
}

type ability struct {
	ID          int64      `json:"id"`
	Multiplier  float64    `json:"multiplier"`
	IsVapeMelt  bool       `json:"is_vape_melt"`
	VapeMeltMul float64    `json:"vape_melt_multiplier"`
	IsPhys      bool       `json:"is_physical"`
	Mods        []modifier `json:"mods"`
}

type artifact struct {
	Level          int64    `json:"level"`
	Slot           slot     `json:"slot"`
	TargetMainStat statType `json:"target_main_stat"`
	MainStat       stat     `json:"main_stat"`
	Substat        []stat   `json:"substat"`
}

type artifactSet struct {
	Set  map[slot]artifact `json:"set"`
	Mods []modifier        `json:"mods"`
}

type stat struct {
	Type  statType `json:"string"`
	Value float64  `json:"value"`
}

type statProb struct {
	Type   statType `json:"type"`
	Weight float64  `json:"weight"`
}

type modifier struct {
	Type  modType `json:"type"`
	Label string  `json:"label"`
	List  []struct {
		Value float64 `json:"value"`
		Desc  string  `json:"desc"`
	} `json:"list"`
}

type slot string
type statType string
type modType string

//Types of artifact slots
const (
	Flower  slot = "Flower"
	Feather slot = "Feather"
	Sands   slot = "Sands"
	Goblet  slot = "Goblet"
	Circlet slot = "Circlet"
)

//stat types
const (
	DEFP statType = "DEF%"
	DEF  statType = "DEF"
	HP   statType = "HP"
	HPP  statType = "HP%"
	ATK  statType = "ATK"
	ATKP statType = "ATK%"
	ER   statType = "ER"
	EM   statType = "EM"
	CR   statType = "CR"
	CD   statType = "CD"
	Heal statType = "Heal"
	EleP statType = "Ele%"
	PhyP statType = "Phys%"
)

//mod types
const (
	ATKPMod  modType = "ATKP_MOD"
	ElePMod  modType = "ELEP_MOD"
	PhyPMod  modType = "PHYP_MOD"
	CRMod    modType = "CR_MOD"
	CDMod    modType = "CD_MOD"
	DmdPMod  modType = "DMGP_MOD"
	EMMod    modType = "EM_MOD"
	ReactMod modType = "REACTION_MOD"
	ResMod   modType = "RESIST_MOD"
	DefMod   modType = "DEF_MOD"
)
