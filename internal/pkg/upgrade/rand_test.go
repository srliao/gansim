package upgrade

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkRand(b *testing.B) {

	s := Simulator{}
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for n := 0; n < b.N; n++ {
		s.rand(rand, ATKP, 20)
	}
}
