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
