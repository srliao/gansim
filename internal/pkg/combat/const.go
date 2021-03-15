package combat

//ElementType is a string representing an element i.e. HYDRO/PYRO/etc...
type ElementType string

//ElementType should be pryo, hydro, cryo, electro, geo, anemo and maybe dendro
const (
	ElementTypePyro    ElementType = "pyro"
	ElementTypeHydro   ElementType = "hydro"
	ElementTypeCryo    ElementType = "cryo"
	ElementTypeElectro ElementType = "electro"
	ElementTypeGeo     ElementType = "geo"
	ElementTypeAnemo   ElementType = "anemo"
)

//ActionType can be swap, dash, jump, attack, skill, burst
type ActionType string

//ActionType constants
const (
	ActionTypeSwap          ActionType = "swap"
	ActionTypeDash          ActionType = "dash"
	ActionTypeJump          ActionType = "jump"
	ActionTypeAttack        ActionType = "attack"
	ActionTypeChargedAttack ActionType = "charge"
	ActionTypeSkill         ActionType = "skill"
	ActionTypeBurst         ActionType = "burst"
)

//CharacterConfig ...
type CharacterConfig struct {
	Name           string             `yaml:"Name"`
	Level          int64              `yaml:"Level"`
	BaseHP         float64            `yaml:"BaseHP"`
	BaseAtk        float64            `yaml:"BaseAtk"`
	BaseDef        float64            `yaml:"BaseDef"`
	BaseCR         float64            `yaml:"BaseCR"`
	BaseCD         float64            `yaml:"BaseCD"`
	AscensionBonus map[string]float64 `yaml:"AscensionBOnus"`
	Constellation  int64              `yaml:"Constellation"`
}

//WeaponConfig ...
type WeaponConfig struct {
	Name          string             `yaml:"Name"`
	BaseAtk       float64            `yaml:"BaseAtk"`
	SecondaryStat map[string]float64 `yaml:"SecondaryStat"`
}

//EnemyConfig ...
type EnemyConfig struct {
	Level      int64   `yaml:"Level"`
	Number     int64   `yaml:"Number"`
	Resistance float64 `yaml:"Resistance"` //this needs to be a map later on
}

//RotationItem ...
type RotationItem struct {
	CharacterName string     `yaml:"CharacterName"`
	Action        ActionType `yaml:"Action"`
	Condition     string     //to be implemented
}

//Config ...
type Config struct {
	Character []CharacterConfig `yaml:"Character"`
	Weapon    WeaponConfig      `yaml:"Weapon"`
	Enemy     EnemyConfig       `yaml:"Enemy"`
	Artifacts map[Slot]Artifact `yaml:"Artifacts"`
	Rotation  []RotationItem
}
