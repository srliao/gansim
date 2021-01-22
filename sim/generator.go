package sim

import (
	"math/rand"
	"time"
)

//Generator for generating random artifacts
type Generator struct {
	rand *rand.Rand
}

//NewGenerator creates a new artifact generator
func NewGenerator(seed int64) *Generator {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	generator := Generator{
		rand: r,
	}

	return &generator

}

//Rand generates one random artifact
func (g *Generator) Rand() {

}

//RandWithMain generates one random artifact with specified main stat
func (g *Generator) RandWithMain() {

}

//RandSet generates one set of random artifact
func (g *Generator) RandSet() {

}

//RandSetWithMain generates set of one random artifact with specified main stat
func (g *Generator) RandSetWithMain() {

}
