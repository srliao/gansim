package combat

import (
	"fmt"
	"math/rand"
	"time"
)

//Aura keeps track of the status of each aura
type Aura struct{}

func (a *Aura) tick() {}

//Field describes field effects (mainly the buffs)
type Field struct {
}

//HookFunc describes a function to be called on a tick; if return true then
//main loop can delete this hook
type HookFunc func(s *Sim) bool

type effectFunc func(s *Sim) bool

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

//Sim keeps track of one simulation
type Sim struct {
	targets    []*Unit
	characters []*character
	active     int
	frame      int

	//per tick hooks
	actions map[string]HookFunc
	//effects
	effects map[effectType]map[string]effectFunc
}

//NewSim creates a new sim unit
func NewSim(n int) *Sim {
	var s Sim
	//setup unit
	var units []*Unit
	for i := 0; i < n; i++ {
		u := &Unit{}
		u.Auras = make(map[ElementType]Aura)
		u.Buffs = make(map[string]int)
		u.Debuffs = make(map[string]int)
		units = append(units, u)
	}
	s.actions = make(map[string]HookFunc)
	s.effects = make(map[effectType]map[string]effectFunc)

	s.targets = units

	return &s
}

func (s *Sim) useEffect(f effectFunc, key string, hook effectType) {
	if _, ok := s.effects[hook]; !ok {
		s.effects[hook] = make(map[string]effectFunc)
	}
	s.effects[hook][key] = f
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
		for _, t := range s.targets {
			t.tick(s)
		}
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

//handleTick
func (s *Sim) handleTick() {
	//apply pre action
	for k, f := range s.effects[preActionHook] {
		if f(s) {
			fmt.Printf("preAction %v expired\n", k)
			delete(s.effects[preActionHook], k)
		}
	}
	//apply actions
	for k, f := range s.effects[actionHook] {
		if f(s) {
			fmt.Printf("action %v expired\n", k)
			delete(s.effects[actionHook], k)
		}
	}
}

//handleAction executes the next action, returns the cooldown
func (s *Sim) handleAction(active int, a Action) int {
	//if active see what ability we want to use
	current := s.characters[active]

	switch a.Type {
	case ActionTypeDash:
		print(s.frame, "dashing")
		return 100
	case ActionTypeJump:
		print(s.frame, "dashing")
		fmt.Printf("[%v] jumping\n", s.frame)
		return 100
	case ActionTypeAttack:
		print(s.frame, "%v executing attack", current.Name)
		return current.attack(s)
	case ActionTypeChargedAttack:
		print(s.frame, "%v executing charged attack", current.Name)
		return current.chargeAttack(s)
	case ActionTypeBurst:
		print(s.frame, "%v executing burst", current.Name)
		return current.burst(s)
	case ActionTypeSkill:
		print(s.frame, "%v executing skill", current.Name)
		return current.skill(s)
	default:
		//do nothing
		print(s.frame, "no action specified: %v. Doing nothing", a.Type)
	}

	return 0
}

func (s *Sim) handleDamage(d DamageProfile) float64 {

	//calculate attack or def
	var a float64
	if d.UseDef {
		a = d.BaseDef*(1+d.Stats[DEFP]) + d.Stats[DEF]
	} else {
		a = d.BaseAtk*(1+d.Stats[ATKP]) + d.Stats[ATK]
	}
	base := d.Multiplier*a + d.FlatDmg

	damage := base * (1 + d.DmgBonus)

	//check if crit
	if rand.Float64() <= d.Stats[CR] {
		damage = damage * (1 + d.Stats[CD])
	}

	//we'll pretend there's only one unit for now...
	for _, u := range s.targets {

		defmod := float64(d.CharacterLevel+100) / (float64(d.CharacterLevel+100) + float64(u.Level+100)*(1-d.DefMod))
		//apply def mod
		damage = damage * defmod
		//apply resist mod
		res := u.Resist + d.ResistMod
		resmod := 1 - res/2
		if res >= 0 && res < 0.75 {
			resmod = 1 - res
		} else if res > 0.75 {
			resmod = 1 / (4*res + 1)
		}
		damage = damage * resmod

		//apply amp mod - TODO
		if d.ApplyAura {

		}

		//apply other multiplier bonus
		if d.OtherMult > 0 {
			damage = damage * d.OtherMult
		}

		u.Damage += damage

		return damage
	}
	return 0
}

//Action describe one action to execute
type Action struct {
	TargetCharIndex int
	Type            ActionType
}
