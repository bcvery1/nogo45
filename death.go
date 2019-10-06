package main

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

var (
	DeathScreen = deathScreen{}
)

type deathScreen struct{}

func (d deathScreen) update(_ float64, win *pixelgl.Window) leveler {
	if win.JustPressed(pixelgl.KeySpace) {
		win.SetClosed(true)
	}

	win.Clear(color.RGBA{
		R: 0x00,
		G: 0x00,
		B: 0x0d,
		A: 0xff,
	})

	return &d
}

// 930008
func (d deathScreen) draw(target pixel.Target) {
	pos := pixel.V(380, 540)
	pos = cam.Unproject(pos)
	t := text.New(pos, atlas)
	t.Color = color.RGBA{
		R: 0x93,
		G: 0x00,
		B: 0x08,
		A: 0xff,
	}
	_, _ = t.WriteString("You Died")
	t.Draw(target, pixel.IM.Scaled(t.Bounds().Center(), 6))

	pos = pixel.V(330, 280)
	pos = cam.Unproject(pos)
	t = text.New(pos, atlas)
	_, _ = t.WriteString("Thanks for playing, press space to exit")
	t.Draw(target, pixel.IM.Scaled(t.Bounds().Center(), 3))

}
