package main

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"

	"github.com/faiface/pixel"
)

var pedestrians [1]*pedestrian

type pedestrian struct {
	pixel.Vec
	X float64
	Y float64
}

func newPedestrian() *pedestrian {
	p := new(pedestrian)
	p.X = 100
	p.Y = 100
	return p
}

func main() {
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	rand.Seed(seed.Int64())

	pedestrians[0] = newPedestrian()
}
