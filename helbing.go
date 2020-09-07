package main

import (
	crand "crypto/rand"
	"image/color"
	"math"
	"math/big"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

const width, height = float64(1024), float64(512)

const (
	ti = 0.5
	t  = 0.5
)

var Ai = math.Pow(2.0, 3)
var Bi = 0.08
var k1 = math.Pow(1.2, 5) // 1.2e5
var k2 = math.Pow(2.4, 5) // 2.4e5

var pedestrians [512]*pedestrian

type pedestrian struct {
	pixel.Vec
	stateCurrent     int
	stateFuture      int
	locationX        [2]float64
	locationY        [2]float64
	desiredVelocityX float64
	desiredVelocityY float64
	desiredVelocity  float64
	velocityX        [2]float64
	velocityY        [2]float64
	weight           float64
	bodyRadius       float64
	C                color.RGBA
}

func newPedestrian() *pedestrian {
	p := new(pedestrian)

	p.stateCurrent = 0
	p.stateFuture = 1

	p.locationX[p.stateCurrent] = random(0, width)
	p.locationY[p.stateCurrent] = random(0, height)

	p.weight = 60
	p.bodyRadius = 0.3

	p.C = color.RGBA{157, 180, 255, 255}

	var pi float64
	if rand.Intn(2) == 0 {
		pi = 0
	} else {
		pi = math.Pi
	}

	var distance float64
	distance = rand.Float64()*0.6 + 1.2

	p.desiredVelocityX = distance * math.Cos(pi)
	p.desiredVelocityY = distance * math.Sin(pi)

	return p
}

func (p *pedestrian) update() {

	if p.locationX[p.stateFuture] > width {
		p.locationX[p.stateFuture] = 0
	} else if p.locationX[p.stateFuture] < 0 {
		p.locationX[p.stateFuture] = width
	}

	if p.locationY[p.stateFuture] > height {
		p.locationY[p.stateFuture] = 0
	} else if p.locationY[p.stateFuture] < 0 {
		p.locationY[p.stateFuture] = height
	}

	tmp := p.stateCurrent
	p.stateCurrent = p.stateFuture
	p.stateFuture = tmp
	//p.X = p.X + 1
	//p.Y = p.Y - 1
}

func (p *pedestrian) draw(imd *imdraw.IMDraw) {
	pix := pixel.V(
		//width/2-p.X,
		//height/2-p.Y,
		width/2-p.locationX[p.stateCurrent],
		height/2-p.locationY[p.stateCurrent],
	)

	fmt.Printf("x: %f, y:%f\n", p.locationX[p.stateCurrent], p.locationY[p.stateCurrent])

	imd.Color = p.C
	imd.Push(pix)
	imd.Circle(5, 0)
}

func main() {
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	rand.Seed(seed.Int64())

	for i := 0; i < len(pedestrians); i++ {
		pedestrians[i] = newPedestrian()
	}

	pixelgl.Run(func() {
		win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
			Bounds:      pixel.R(0, 0, width, height),
			VSync:       true,
			Undecorated: false,
		})
		if err != nil {
			panic(err)
		}

		imd := imdraw.New(nil)
		imd.Precision = 7
		imd.SetMatrix(pixel.IM.Moved(win.Bounds().Center()))

		for !win.Closed() {
			win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))

			imd.Clear()

			for _, s := range pedestrians {
				s.update()
				s.draw(imd)
			}

			win.Clear(color.Black)
			imd.Draw(win)
			win.Update()
		}
	},
	)
}

func random(min, max float64) float64 {
	return rand.Float64() * (max - min)
}
