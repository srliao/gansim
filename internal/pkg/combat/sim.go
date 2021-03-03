package combat

//Unit keeps track of the status of one enemy Unit
type Unit struct {
	Level float64

	//Auras
	Auras   map[Elements]Aura
	Buffs   map[string]bool
	Debuffs map[string]bool

	//Hooks are called each tick; not sure usage yet
	Hooks map[string]func()

	//special hooks
	DamageHooks map[string]HookFunc //this is for actions that triggers on damage, such as fischl oz thinggy

	//stats
	Damage float64 //total damage received
}

//ApplyDamage applies ability damage to an Unit
func (u *Unit) ApplyDamage() {}

//ApplyAura applies an aura to the Unit
func (u *Unit) ApplyAura() {
	//can trigger apply damage for superconduct, electrocharged, etc..
}

func (u *Unit) tick() {}

//Field keeps track of field statuses
type Field struct {
}

//Aura keeps track of the status of each aura
type Aura struct{}

func (a *Aura) tick() {}

//Elements is a string representing an element i.e. HYDRO/PYRO/etc...
type Elements string

//Action can be swap, dash, jump, attack, skill, burst
type Action string

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
	OnField map[Action]int

	//ICD tracker for actions
	ICD map[Action]int

	//TickHooks are functions to be called on each tick
	//this is useful for on field effect such as gouba/oz/pyronado
	//we can use store to keep track of the uptime on gouba/oz/pyronado/taunt etc..
	//for something like baron bunny, if uptime = xx, then trigger damage
	TickHooks map[string]HookFunc
	//what about something like bennett ult or ganyu ult that affects char in the field?? this hook can only affect current actor?

	//ability functions to be defined by each character on how they will
	//affect the unit
	Attack func(u *Unit)
	Skill  func(u *Unit)
	Burst  func(u *Unit)

	//somehow we have to deal with artifact effects too?
	ArtifactSetBonus func(u *Unit)

	//character specific information; need this for damage calc
	Level         int64
	BaseAtk       int64
	WeaponAtk     int64
	ArtifactStats map[string]float64
	Talent        map[Action]int64 //talent levels
}

//HookFunc describes a function to be called on a tick
type HookFunc func(u *Unit)

func (c *Character) tick(u *Unit) {
	//this function gets called for every character every tick
}

func (c *Character) orb(e Elements, isActive bool) {
	//called when elemental orgs are received by the character
}

//Sim keeps track of one simulation
type Sim struct {
	Target *Unit
	Actors []Character
}

//Run the sim
func (s *Sim) Run() {
	var cooldown int
	var active int //index of the currently active car
	//60fps, 60s/min, 2min
	for f := 0; f < 7200; f++ {
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
		next := s.next()

		//we trigger the next action; action currently is either a character abil
		//or a swap
		//this corresponds to a manual action in game by the player i.e.
		//left click, right click, q, or e press
		switch next {
		case "swap":
			active = 0     //swap active to whoever
			cooldown = 100 //trigger cooldown
		case "dash":
			//simply puts us in cooldown
			cooldown = 100
		case "jump":
			//simply puts us in cooldown
			cooldown = 100
		case "attack":
			s.Actors[active].Attack(s.Target)
		case "burst":
		case "skill":
		}
	}
}

func (s *Sim) next() Action {
	//determine the next action somehow
	return "swap"
}
