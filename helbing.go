package main

import (
	crand "crypto/rand"
	"fmt"
	"image/color"
	"math"
	"math/big"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

const width, height = float64(8), float64(4)
const ratio = float64(100)

const (
	ti = 0.03
	t  = 0.03
)

// Ai is a constant of repulsive interaction force.
var Ai = math.Pow(2.0, 3)

// Ai is a constant of repulsive interaction force.
var Bi = 0.08

// k1 is a constant of body force.
var k1 = math.Pow(1.2, 5) // 1.2e5

// k2 is a constant of sliding friction force.
var k2 = math.Pow(2.4, 5) // 2.4e5

var pedestrians [256]*pedestrian

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

	var pi float64
	if rand.Intn(2) == 0 {
		pi = 0
		p.C = color.RGBA{0, 255, 0, 255}
	} else {
		pi = math.Pi
		p.C = color.RGBA{255, 0, 255, 255}
	}

	var distance float64
	distance = rand.Float64()*0.6 + 1.2

	p.desiredVelocityX = distance * math.Cos(pi)
	p.desiredVelocityY = distance * math.Sin(pi)

	return p
}

func (p *pedestrian) update() {

	if p.locationX[p.stateFuture] > width {
		p.locationX[p.stateFuture] = p.locationX[p.stateFuture] - width
	} else if p.locationX[p.stateFuture] < 0 {
		p.locationX[p.stateFuture] = width + p.locationX[p.stateFuture]
	}

	if p.locationY[p.stateFuture] > height {
		p.locationY[p.stateFuture] = p.locationY[p.stateFuture] - height
	} else if p.locationY[p.stateFuture] < 0 {
		p.locationY[p.stateFuture] = height + p.locationY[p.stateFuture]
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
		(width/2-p.locationX[p.stateCurrent])*ratio,
		(height/2-p.locationY[p.stateCurrent])*ratio,
	)

	fmt.Printf("x: %f, y:%f\n", p.locationX[p.stateCurrent], p.locationY[p.stateCurrent])

	imd.Color = p.C
	imd.Push(pix)
	imd.Circle(p.bodyRadius*ratio/5, 1)
}

func main() {
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	rand.Seed(seed.Int64())

	for i := 0; i < len(pedestrians); i++ {
		pedestrians[i] = newPedestrian()
	}

	pixelgl.Run(func() {
		win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
			Bounds:      pixel.R(0, 0, width*ratio, height*ratio),
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

			for i := 0; i < len(pedestrians); i++ {

				pi := pedestrians[i]

				m := pi.weight

				var fX, fY float64

				maX := m * (pi.desiredVelocityX - pi.velocityX[pi.stateCurrent]) / ti
				maY := m * (pi.desiredVelocityY - pi.velocityY[pi.stateCurrent]) / ti

				fX += maX
				fY += maY

				var fijX, fijY float64

				for j := 0; j < len(pedestrians); j++ {
					if i == j {
						continue
					}
					pj := pedestrians[j]

					diffX := pi.locationX[pi.stateCurrent] - pj.locationX[pj.stateCurrent]
					diffY := pi.locationY[pi.stateCurrent] - pj.locationY[pj.stateCurrent]

					disX := math.Abs(diffX)
					disY := math.Abs(diffY)

					if disX >= -15 && disX <= 15 && disY >= -15 && disY <= 15 {

						dij := math.Sqrt(disX*disX + disY*disY)

						nijX := diffX / dij
						nijY := diffY / dij

						rij := pi.bodyRadius + pj.bodyRadius

						tijX := -nijY
						tijY := nijX

						dvjitX := (pj.velocityX[pj.stateCurrent] - pi.velocityX[pi.stateCurrent]) * tijX
						dvjitY := (pj.velocityY[pj.stateCurrent] - pi.velocityY[pi.stateCurrent]) * tijY

						funcG := rij - dij
						if funcG < 0 {
							funcG = 0
						}

						appr := (rij - dij) / Bi

						fExp := Ai*math.Exp(appr) + k1*funcG

						fijX += fExp*nijX + k2*funcG*dvjitX*tijX
						fijY += fExp*nijY + k2*funcG*dvjitY*tijY
					}
				}

				fX += fijX
				fY += fijY

				v := math.Sqrt(math.Pow(pi.desiredVelocityX, 2)+math.Pow(pi.desiredVelocityY, 2)) / t * m

				fLength := math.Sqrt(math.Pow(fX, 2) + math.Pow(fY, 2))

				if fX > v {
					fX = fX / fLength * v
				}
				if fX < -v {
					fX = fX / fLength * v
				}
				if fY > v {
					fY = fY / fLength * v
				}
				if fY < -v {
					fY = fY / fLength * v
				}

				pi.velocityX[pi.stateFuture] = (fX / m) * t
				pi.velocityY[pi.stateFuture] = (fY / m) * t

				pi.locationX[pi.stateFuture] = pi.locationX[pi.stateCurrent] + (fX/(2*m))*t*t
				pi.locationY[pi.stateFuture] = pi.locationY[pi.stateCurrent] + (fY/(2*m))*t*t

				pi.draw(imd)
				pi.update()
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
