package main

import (
	"image/color"
	"strconv"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

var (
	HUD = &hud{}

	coinPic *pixel.Sprite
)

type hud struct{}

func (h *hud) update(dt float64, win *pixelgl.Window) leveler {
	return currentLvl
}

func (h *hud) draw(target pixel.Target) {
	pos := pixel.V(winBounds.Max.X-32, 32)
	unproj := cam.Unproject(pos)
	coinPic.Draw(target, pixel.IM.Moved(unproj).Scaled(unproj, 1.5))

	t := text.New(unproj, atlas)
	t.Color = defaultTextColour
	_, _ = t.WriteString(strconv.Itoa(Player.coins))
	t.Draw(target, pixel.IM.Moved(pixel.V(-4, -24)))

	if Player.health < 100 && Player.health > 0 {
		imd := imdraw.New(nil)

		// border
		bottomL := cam.Unproject(pixel.V(19, winBounds.Max.Y-20-101))
		topR := cam.Unproject(pixel.V(20+11, winBounds.Max.Y-19))

		imd.Color = color.RGBA{
			R: 0x00,
			G: 0x00,
			B: 0x00,
			A: 0xff,
		}
		imd.Push(bottomL, topR)
		imd.Rectangle(2)

		// bar
		bottomL = cam.Unproject(pixel.V(20, winBounds.Max.Y-20-100))
		topR = cam.Unproject(pixel.V(20+10, winBounds.Max.Y-20-(100-Player.health)))

		imd.Color = color.RGBA{
			R: 0xcd,
			G: 0x00,
			B: 0x0b,
			A: 0xff,
		}
		imd.Push(bottomL, topR)
		imd.Rectangle(0)

		imd.Draw(target)
	}
}
