package combat

//character contains all the information required to calculate
type character struct {
	//track cooldowns in general; can be skill on field, ICD, etc...
	cooldown map[string]int
	energy   float64 //how much energy the character currently have
	stamina  float64 //how much stam the character currently have

	//we need some sort of key/val store to store information
	//specific to each character.
	//use to keep track of attack counter/diluc e counter, etc...
	store map[string]interface{}

	//Init is used to add in any initial hooks to the sim
	init func(s *Sim)

	//tickHooks are functions to be called on each tick
	//this is useful for on field effect such as gouba/oz/pyronado
	//we can use store to keep track of the uptime on gouba/oz/pyronado/taunt etc..
	//for something like baron bunny, if uptime = xx, then trigger damage
	tickHooks map[string]HookFunc
	//what about something like bennett ult or ganyu ult that affects char in the field?? this hook can only affect current actor?

	//ability functions to be defined by each character on how they will
	//affect the unit
	attack       func(s *Sim) int
	chargeAttack func(s *Sim) int
	plungeAttack func(s *Sim) int
	skill        func(s *Sim) int
	burst        func(s *Sim) int

	//somehow we have to deal with artifact effects too?
	ArtifactSetBonus func(u *Unit)

	//key stats
	stats    map[StatType]float64
	statMods map[string]map[StatType]float64 //special effect mods (character only)

	//character specific information; need this for damage calc
	Name      string
	Level     int64
	BaseAtk   float64
	WeaponAtk float64
	BaseDef   float64
	Talent    map[ActionType]int64 //talent levels
}

func (c *character) tick(s *Sim) {
	//this function gets called for every character every tick
	for k, v := range c.cooldown {
		if v == 0 {
			delete(c.cooldown, k)
		} else {
			c.cooldown[k]--
		}
	}
}

func (c *character) orb(e ElementType, isActive bool) {
	//called when elemental orgs are received by the character
}
