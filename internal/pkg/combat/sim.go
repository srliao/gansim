package combat

import (
	"fmt"
	"math/rand"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type AbilFunc func(s *Sim) int
type ActionFunc func(s *Sim) bool

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
	Target     *unit
	Characters []*Character
	Active     int
	Frame      int

	//per tick hooks
	actions map[string]ActionFunc
	//effects
	effects map[effectType]map[string]effectFunc
}

//New creates new sim from given profile
func New(p Profile) (*Sim, error) {
	s := &Sim{}

	u := &unit{}

	u.auras = make(map[eleType]aura)
	u.status = make(map[string]int)
	u.Level = p.Enemy.Level
	u.Resist = p.Enemy.Resist

	s.Target = u

	s.actions = make(map[string]ActionFunc)
	s.effects = make(map[effectType]map[string]effectFunc)

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	switch p.LogLevel {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	}
	config.EncoderConfig.TimeKey = ""

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(logger)

	var chars []*Character
	//create the characters
	for _, v := range p.Characters {
		//initialize artifact sets

		f, ok := charMap[v.Name]
		if !ok {
			return nil, fmt.Errorf("invalid character: %v", v.Name)
		}

		c := f(s, logger.Sugar())
		//initialize other variables/stats
		c.Stats = make(map[StatType]float64)
		c.Cooldown = make(map[string]int)
		c.Store = make(map[string]interface{})
		c.Mods = make(map[string]map[StatType]float64)
		c.Profile = v

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
			c.Stats[a.MainStat.Type] += a.MainStat.Value
			for _, sub := range a.Substat {
				c.Stats[sub.Type] += sub.Value
			}
			sb[a.Set]++
		}
		//add ascension bonus
		for k, v := range v.AscensionBonus {
			c.Stats[k] += v
		}
		//add weapon sub stat
		for k, v := range v.WeaponSecondaryStat {
			c.Stats[k] += v
		}
		//add set bonus
		for key, count := range sb {
			if f, ok := setBonus[key]; ok {
				f(c, s, count)
			}
		}

		chars = append(chars, c)
	}
	s.Characters = chars

	return s, nil
}

//Run the sim; length in seconds
func (s *Sim) Run(length int, list []Action) float64 {
	var cooldown int
	var active int //index of the currently active car
	var i int
	rand.Seed(time.Now().UnixNano())
	//60fps, 60s/min, 2min
	for s.Frame = 0; s.Frame < 60*length; s.Frame++ {
		//tick target and each character
		//target doesn't do anything, just takes punishment, so it won't affect cd
		s.Target.tick(s)
		for _, c := range s.Characters {
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
			fmt.Printf("[%v] swapping to char #%v (current = %v)\n", s.Frame, next.TargetCharIndex, active)
			//trigger a swap
			cooldown = 150
			active = next.TargetCharIndex
			continue

		}
		//move on to next action on list
		i++

		cooldown = s.handleAction(active, next)

	}

	return s.Target.damage
}

func (s *Sim) addEffect(f effectFunc, key string, hook effectType) {
	if _, ok := s.effects[hook]; !ok {
		s.effects[hook] = make(map[string]effectFunc)
	}
	s.effects[hook][key] = f
}

func (s *Sim) AddAction(f ActionFunc, key string) {
	s.actions[key] = f
}

//handleTick
func (s *Sim) handleTick() {
	for k, f := range s.actions {
		if f(s) {
			print(s.Frame, true, "action %v expired", k)
			delete(s.actions, k)
		}
	}
}

//handleAction executes the next action, returns the cooldown
func (s *Sim) handleAction(active int, a Action) int {
	//if active see what ability we want to use
	c := s.Characters[active]

	switch a.Type {
	case ActionTypeDash:
		print(s.Frame, false, "dashing")
		return 100
	case ActionTypeJump:
		print(s.Frame, false, "dashing")
		fmt.Printf("[%v] jumping\n", s.Frame)
		return 100
	case ActionTypeAttack:
		print(s.Frame, false, "%v executing attack", c.Profile.Name)
		return c.Attack(s)
	case ActionTypeChargedAttack:
		print(s.Frame, false, "%v executing charged attack", c.Profile.Name)
		return c.ChargeAttack(s)
	case ActionTypeBurst:
		print(s.Frame, false, "%v executing burst", c.Profile.Name)
		return c.Burst(s)
	case ActionTypeSkill:
		print(s.Frame, false, "%v executing skill", c.Profile.Name)
		return c.Skill(s)
	default:
		//do nothing
		print(s.Frame, false, "no action specified: %v. Doing nothing", a.Type)
	}

	return 0
}

//Action describe one action to execute
type Action struct {
	TargetCharIndex int
	Type            ActionType
}

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

type Profile struct {
	Label      string             `yaml:"Label"`
	Enemy      EnemyProfile       `yaml:"Enemy"`
	Characters []CharacterProfile `yaml:"Characters"`
	Rotation   []RotationItem     `yaml:"Rotation"`
	LogLevel   string             `yaml:"LogLevel"`
}

//CharacterProfile ...
type CharacterProfile struct {
	Name                string               `yaml:"Name"`
	Level               int64                `yaml:"Level"`
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
}

//EnemyProfile ...
type EnemyProfile struct {
	Level  int64   `yaml:"Level"`
	Resist float64 `yaml:"Resistance"` //this needs to be a map later on
}

//RotationItem ...
type RotationItem struct {
	CharacterName string     `yaml:"CharacterName"`
	Action        ActionType `yaml:"Action"`
	Condition     string     //to be implemented
}
