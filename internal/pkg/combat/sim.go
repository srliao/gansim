package combat

import (
	"fmt"
	"math/rand"
	"time"
)

//Aura keeps track of the status of each aura
type Aura struct{}

func (a *Aura) tick() {}

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

//Character contains all the information required to calculate
type Character struct {
	//track cd of the abilities
	Cooldown map[string]int
	Energy   float64 //how much energy the character currently have

	//info specific to normal attacks; SHOULD be identical for all char
	//hope mihoyo doesnt come out with some weird ass char that changes this
	AttackCounter    int //which attack in the chain are we in; starts at 0 (= 1)
	AttackResetTimer int //number of frames before attack counter resets back to 0

	//we need some sort of key/val Store to Store information
	//specific to each character. not sure what this could be yet
	Store map[string]interface{}

	//OnField uptime tracker for any skill that stays on field independent
	//of user action
	OnField map[ActionType]int

	//ICD tracker for actions
	ICD map[ActionType]int

	//TickHooks are functions to be called on each tick
	//this is useful for on field effect such as gouba/oz/pyronado
	//we can use store to keep track of the uptime on gouba/oz/pyronado/taunt etc..
	//for something like baron bunny, if uptime = xx, then trigger damage
	TickHooks map[string]HookFunc
	//what about something like bennett ult or ganyu ult that affects char in the field?? this hook can only affect current actor?

	//ability functions to be defined by each character on how they will
	//affect the unit
	Attack       func(s *Sim) int
	ChargeAttack func(s *Sim) int
	Skill        func(s *Sim) int
	Burst        func(s *Sim) int

	//somehow we have to deal with artifact effects too?
	ArtifactSetBonus func(u *Unit)

	//character specific information; need this for damage calc
	Level         int64
	BaseAtk       float64
	WeaponAtk     float64
	BaseDef       float64
	ArtifactStats map[StatType]float64
	Talent        map[ActionType]int64 //talent levels
}

//HookFunc describes a function to be called on a tick
type HookFunc func(u *Unit)

func (c *Character) tick(u *Unit) {
	//this function gets called for every character every tick
}

func (c *Character) orb(e ElementType, isActive bool) {
	//called when elemental orgs are received by the character
}

//Field describes field effects (mainly the buffs)
type Field struct {
}

//Sim keeps track of one simulation
type Sim struct {
	Target *Unit
	Actors []*Character
	Field  *Field
	Active int
	Frame  int
}

//Run the sim; length in seconds
func (s *Sim) Run(length int, list []Action) {
	var cooldown int
	var active int //index of the currently active car
	var i int
	rand.Seed(time.Now().UnixNano())
	//60fps, 60s/min, 2min
	for s.Frame = 0; s.Frame < 60*length; s.Frame++ {
		//tick target and each character
		//target doesn't do anything, just takes punishment, so it won't affect cd
		s.Target.tick()
		for _, c := range s.Actors {
			//character may affect cooldown by i.e. adding to it
			//character also need to know if we're currently on cooldown
			//so they don't do anything other than tick down
			c.tick(s.Target)
		}

		//if in cooldown, do nothing
		if cooldown > 0 {
			cooldown--
			continue
		}

		//otherwise only either action or swaps can trigger cooldown
		//we figure out what the next action is to be
		next := list[i]

		//check if actor is active
		if next.TargetCharIndex != active {
			fmt.Printf("[%v] swapping to char #%v (current = %v)\n", s.Frame, next.TargetCharIndex, active)
			//trigger a swap
			cooldown = 150
			active = next.TargetCharIndex
			continue

		}

		//if active see what ability we want to use
		current := s.Actors[active]
		switch next.Type {
		case ActionTypeDash:
			fmt.Printf("[%v] dashing\n", s.Frame)
			cooldown = 100
		case ActionTypeJump:
			fmt.Printf("[%v] jumping\n", s.Frame)
			cooldown = 100
		case ActionTypeAttack:
			fmt.Printf("[%v] char #%v executing attack\n", s.Frame, active)
			cooldown = current.Attack(s)
		case ActionTypeChargedAttack:
			fmt.Printf("[%v] char #%v executing charged attack\n", s.Frame, active)
			cooldown = current.ChargeAttack(s)
		case ActionTypeBurst:
			fmt.Printf("[%v] char #%v executing burst\n", s.Frame, active)
			cooldown = current.Burst(s)
		case ActionTypeSkill:
			fmt.Printf("[%v] char #%v executing skill\n", s.Frame, active)
			cooldown = current.Skill(s)
		default:
			//do nothing
			fmt.Printf("[%v] no action specified: %v. Doing nothing\n", s.Frame, next.Type)
		}
		//move on to next action on list
		i++
	}
}

func (s *Sim) next() ActionType {
	//determine the next action somehow
	return "swap"
}

//Action describe one action to execute
type Action struct {
	TargetCharIndex int
	Type            ActionType
}
