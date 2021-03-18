package combat

import (
	"fmt"
	"math/rand"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Profile struct {
	Label string `yaml:"Label"`
	Enemy struct {
		Level int64 `yaml:"Level"`
		// Number     int64   `yaml:"Number"`
		Resist float64 `yaml:"Resist"`
	} `yaml:"Enemy"`
	Characters []struct {
		Name                string               `yaml:"Name"`
		Level               int                  `yaml:"Level"`
		BaseHP              float64              `yaml:"BaseHP"`
		BaseAtk             float64              `yaml:"BaseAtk"`
		BaseDef             float64              `yaml:"BaseDef"`
		BaseCR              float64              `yaml:"BaseCR"`
		BaseCD              float64              `yaml:"BaseCD"`
		Constellation       int                  `yaml:"Constellation"`
		AscensionBonus      map[StatType]float64 `yaml:"AscensionBonus"`
		WeaponName          string               `yaml:"WeaponName"`
		WeaponRefinement    int                  `yaml:"WeaponRefinement"`
		WeaponBaseAtk       float64              `yaml:"WeaponBaseAtk"`
		WeaponSecondaryStat map[StatType]float64 `yaml:"WeaponSecondaryStat"`
		Artifacts           map[Slot]Artifact    `yaml:"Artifacts"`
	} `yaml:"Characters"`
	Rotation  []RotationItem `yaml:"Rotation"`
	ShowDebug bool           `yaml:"ShowDebug"`
}

//RotationItem ...
type RotationItem struct {
	CharacterName string     `yaml:"CharacterName"`
	Action        ActionType `yaml:"Action"`
	Condition     string     //to be implemented
}

type actionFunc func(s *Sim) bool

type effectType string

const (
	preDamageHook   effectType = "PRE_DAMAGE"
	postDamageHook  effectType = "POST_DAMAGE"
	preAuraAppHook  effectType = "PRE_AURA_APP"
	postAuraAppHook effectType = "POST_AURA_APP"
	preActionHook   effectType = "PRE_ACTION"
	actionHook      effectType = "ACTION"
	postActionHook  effectType = "POST_ACTION"
	fieldEffectHook effectType = "FIELD_EFFECt"
)

type effectFunc func(s *snapshot) bool

//Sim keeps track of one simulation
type Sim struct {
	target     *unit
	characters []*character
	active     int
	frame      int

	//per tick hooks
	actions map[string]actionFunc
	//effects
	effects map[effectType]map[string]effectFunc
}

//New creates new sim from given profile
func New(p Profile) (*Sim, error) {
	s := &Sim{}

	u := &unit{}

	u.auras = make(map[eleType]int)
	u.status = make(map[string]int)
	u.Level = p.Enemy.Level
	u.Resist = p.Enemy.Resist

	s.target = u

	s.actions = make(map[string]actionFunc)
	s.effects = make(map[effectType]map[string]effectFunc)
	var chars []*character
	//create the characters
	for _, v := range p.Characters {
		c := &character{}
		//initialize artifact sets
		//initialize other variables/stats
		c.stats = make(map[StatType]float64)
		c.cooldown = make(map[string]int)
		c.store = make(map[string]interface{})
		c.statMods = make(map[string]map[StatType]float64)
		c.BaseAtk = v.BaseAtk
		c.BaseDef = v.BaseDef
		c.BaseHP = v.BaseHP
		c.BaseCD = v.BaseCD
		c.BaseCR = v.BaseCR

		switch v.Name {
		case "Ganyu":
			newGanyu(c)
		default:
			return nil, fmt.Errorf("invalid character: %v", v.Name)
		}
		//initialize weapon
		switch v.WeaponName {
		case "Prototype Crescent":
			weaponPrototypeCrescent(c, s, v.WeaponRefinement)
		default:
			return nil, fmt.Errorf("invalid weapon: %v", v.WeaponName)
		}
		c.WeaponAtk = v.WeaponBaseAtk
		//check set bonus
		sb := make(map[string]int)
		for _, a := range v.Artifacts {
			c.stats[a.MainStat.Type] += a.MainStat.Value
			for _, sub := range a.Substat {
				c.stats[sub.Type] += sub.Value
			}
			sb[a.Set]++
		}
		//add ascension bonus
		for k, v := range v.AscensionBonus {
			c.stats[k] += v
		}
		//add weapon sub stat
		for k, v := range v.WeaponSecondaryStat {
			c.stats[k] += v
		}
		//add set bonus
		for key, count := range sb {
			if f, ok := setBonus[key]; ok {
				f(c, s, count)
			}
		}

		chars = append(chars, c)
	}
	s.characters = chars

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	if !p.ShowDebug {
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
	config.EncoderConfig.TimeKey = ""

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(logger)

	return s, nil
}

//Run the sim; length in seconds
func (s *Sim) Run(length int, list []Action) {
	var cooldown int
	var active int //index of the currently active car
	var i int
	rand.Seed(time.Now().UnixNano())
	//60fps, 60s/min, 2min
	for s.frame = 0; s.frame < 60*length; s.frame++ {
		//tick target and each character
		//target doesn't do anything, just takes punishment, so it won't affect cd
		s.target.tick(s)
		for _, c := range s.characters {
			//character may affect cooldown by i.e. adding to it
			c.tick(s)
		}

		s.handleTick()

		//if in cooldown, do nothing
		if cooldown > 0 {
			cooldown--
			continue
		}

		if i >= len(list) {
			//start over
			i = 0
		}
		//otherwise only either action or swaps can trigger cooldown
		//we figure out what the next action is to be
		next := list[i]

		//check if actor is active
		if next.TargetCharIndex != active {
			fmt.Printf("[%v] swapping to char #%v (current = %v)\n", s.frame, next.TargetCharIndex, active)
			//trigger a swap
			cooldown = 150
			active = next.TargetCharIndex
			continue

		}
		//move on to next action on list
		i++

		cooldown = s.handleAction(active, next)

	}
}

func (s *Sim) addEffect(f effectFunc, key string, hook effectType) {
	if _, ok := s.effects[hook]; !ok {
		s.effects[hook] = make(map[string]effectFunc)
	}
	s.effects[hook][key] = f
}

func (s *Sim) addAction(f actionFunc, key string) {
	s.actions[key] = f
}

//handleTick
func (s *Sim) handleTick() {
	for k, f := range s.actions {
		if f(s) {
			print(s.frame, true, "action %v expired", k)
			delete(s.actions, k)
		}
	}
}

//handleAction executes the next action, returns the cooldown
func (s *Sim) handleAction(active int, a Action) int {
	//if active see what ability we want to use
	current := s.characters[active]

	switch a.Type {
	case ActionTypeDash:
		print(s.frame, false, "dashing")
		return 100
	case ActionTypeJump:
		print(s.frame, false, "dashing")
		fmt.Printf("[%v] jumping\n", s.frame)
		return 100
	case ActionTypeAttack:
		print(s.frame, false, "%v executing attack", current.Name)
		return current.attack(s)
	case ActionTypeChargedAttack:
		print(s.frame, false, "%v executing charged attack", current.Name)
		return current.chargeAttack(s)
	case ActionTypeBurst:
		print(s.frame, false, "%v executing burst", current.Name)
		return current.burst(s)
	case ActionTypeSkill:
		print(s.frame, false, "%v executing skill", current.Name)
		return current.skill(s)
	default:
		//do nothing
		print(s.frame, false, "no action specified: %v. Doing nothing", a.Type)
	}

	return 0
}

//Action describe one action to execute
type Action struct {
	TargetCharIndex int
	Type            ActionType
}
