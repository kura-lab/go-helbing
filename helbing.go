package main

import (
	"github.com/faiface/pixel"
)

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
