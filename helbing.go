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

var pedestrians [1]*pedestrian

type pedestrian struct {
	pixel.Vec
	X float64
	Y float64
	C color.RGBA
}

func newPedestrian() *pedestrian {
	p := new(pedestrian)
	p.X = 100
	p.Y = 100
	p.C = color.RGBA{157, 180, 255, 255}
	return p
}

func (p *pedestrian) update() {

	if p.X > 200 {
		p.X = 100
		p.Y = 100
	}

	p.X = p.X + 1
	p.Y = p.Y - 1
}

func (p *pedestrian) draw(imd *imdraw.IMDraw) {
	pix := pixel.V(
		p.X,
		p.Y,
	)

	imd.Color = p.C
	imd.Push(pix)
	imd.Circle(5, 0)
}

func main() {
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	rand.Seed(seed.Int64())

	pedestrians[0] = newPedestrian()

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
