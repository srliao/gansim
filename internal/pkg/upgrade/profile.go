package upgrade

type Profile struct {
	Label         string  `yaml:"Label"`
	CharLevel     float64 `yaml:"CharacterLevel"`
	CharBaseAtk   float64 `yaml:"CharacterBaseAtk"`
	WeaponBaseAtk float64 `yaml:"WeaponBaseAtk"`
	Enemy         struct {
		Level  float64             `yaml:"Level"`
		Resist map[EleType]float64 `yaml:"Resist"`
	}
	//artifact details
	Artifacts   []Artifact `yaml:"CurrentArtifacts"`
	NewArtifact struct {
		DesiredLevel int      `yaml:"DesiredLevel"`
		Artifact     Artifact `yaml"Artifact"`
	} `yaml:"New"`
	GlobalStatMod map[string][]float64 `yaml:"GlobalStatMod"`
	//abilities
	Abilities []struct {
		DamageMult        float64               `default:"1.0" yaml:"DamageMult"`
		Element           EleType               `default:"physical" yaml:"Element"`
		VapMeltMultiplier float64               `default:"1.0" yaml:"VaporizeOrMeltMultiplier"`
		StatMods          map[string][]float64  `yaml:"StatMods"`
		ReactionBonus     []float64             `yaml:"ReactionBonus"`
		ResistMod         map[EleType][]float64 `yaml:"ResistMod"`
		DefShredMod       []float64             `yaml:"DefShredMod"`
	} `yaml:"Abilities"`
}

type Artifact struct {
	Level int                `yaml:"Level"`
	Set   string             `yaml:"Set"`
	Main  map[string]float64 `yaml:"Main"`
	Sub   map[string]float64 `yaml:"Sub"`
}

//EleType is a string representing an element i.e. HYDRO/PYRO/etc...
type EleType string

//ElementType should be pryo, Hydro, Cryo, Electro, Geo, Anemo and maybe dendro
const (
	Pyro      EleType = "pyro"
	Hydro     EleType = "hydro"
	Cryo      EleType = "cryo"
	Electro   EleType = "electro"
	Geo       EleType = "geo"
	Anemo     EleType = "anemo"
	Dendro    EleType = "dendro"
	Physical  EleType = "physical"
	NoElement EleType = ""
)
