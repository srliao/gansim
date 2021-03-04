package combat

import (
	"testing"
)

func TestSim(t *testing.T) {

	g := NewGanyu()
	var s Sim
	enemy := Unit{
		Level:  100,
		Resist: 0.1,
	}
	s.Target = &enemy
	s.Actors = append(s.Actors, g)

	s.Active = 0
	var actions = []Action{
		{
			TargetCharIndex: 0,
			Type:            ActionTypeChargedAttack,
		},
		{
			TargetCharIndex: 0,
			Type:            ActionTypeChargedAttack,
		},
	}
	s.Run(60, actions)

}
