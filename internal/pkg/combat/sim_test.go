package combat

import (
	"io/ioutil"
	"log"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestSim(t *testing.T) {

	var source []byte
	var cfg Profile
	var err error

	source, err = ioutil.ReadFile("./test/cfg.yaml")
	if err != nil {
		t.Fatal(err)
	}
	err = yaml.Unmarshal(source, &cfg)
	if err != nil {
		t.Fatal(err)
	}

	s, err := New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	var actions = []Action{
		// {
		// 	TargetCharIndex: 0,
		// 	Type:            ActionTypeBurst,
		// },
		{
			TargetCharIndex: 0,
			Type:            ActionTypeChargedAttack,
		},
	}
	s.Run(6, actions)

	// s := New()

	// g := newGanyu()
	// g.BaseAtk = cfg.Profile.CharBaseAtk
	// g.Level = cfg.Profile.CharLevel
	// g.WeaponAtk = cfg.Profile.WeaponBaseAtk

	// for _, a := range cfg.Artifacts {
	// 	g.stats[a.MainStat.Type] += a.MainStat.Value
	// 	for _, v := range a.Substat {
	// 		g.stats[v.Type] += v.Value
	// 	}
	// }
	// //manually add weapon mods etc..
	// g.stats[ATKP] += 0.413
	// g.stats[CryoP] += 0.15
	// g.stats[CR] += 0.25
	// g.stats[CD] += 0.884

	// s.characters = append(s.characters, g)

	// s.active = 0

}

type TestProfile struct {
	Label         string  `yaml:"Label"`
	CharLevel     int64   `yaml:"CharacterLevel"`
	CharBaseAtk   float64 `yaml:"CharacterBaseAtk"`
	WeaponBaseAtk float64 `yaml:"WeaponBaseAtk"`
	EnemyLevel    int64   `yaml:"EnemyLevel"`
	//artifact details
	Artifacts struct {
		Level          int64             `default:"20" yaml:"Level"`
		Current        map[Slot]Artifact `yaml:"Current"`
		TargetMainStat map[Slot]StatType `yaml:"TargetMainStat"`
	} `yaml:"Artifacts"`
	//abilities
	Abilities []struct {
		Talent            float64   `default:"1.0" yaml:"Talent"`
		TalentIsElemental bool      `default:"true" yaml:"TalentIsElemental"`
		VapMeltMultiplier float64   `default:"1.0" yaml:"VaporizeOrMeltMultiplier"`
		AtkMod            []float64 `yaml:"AtkMod"`
		EleMod            []float64 `yaml:"EleMod"`
		PhyMod            []float64 `yaml:"PhyMod"`
		CCMod             []float64 `yaml:"CCMod"`
		CDMod             []float64 `yaml:"CDMod"`
		DmgMod            []float64 `yaml:"DmgMod"`
		EMMod             []float64 `yaml:"EMMod"`
		ReactionBonus     []float64 `yaml:"ReactionBonus"`
		ResistMod         []float64 `yaml:"ResistMod"`
		DefShredMod       []float64 `yaml:"DefShredMod"`
		//special stat modifiers such as staff of honma
		SpecialStatMod []struct {
			BonusStat StatType `yaml:"BonusStat"`
			ScaleStat StatType `yaml:"ScaleStat"`
			Modifier  float64  `yaml:"Modifier"`
		} `yaml:"SpecialStatMod"`
		//special dmg modifiers such as zhongli
		SpecialDmgMod []struct {
			ScaleStat StatType `yaml:"ScaleStat"`
			Modifier  float64  `yaml:"Modifier"`
		} `yaml:"SpecialDmgMod"`
	} `yaml:"Abilities"`
}

type testCfg struct {
	Profile   TestProfile       `yaml:"Profile"`
	Artifacts map[Slot]Artifact `yaml:"Artifacts"`
}
