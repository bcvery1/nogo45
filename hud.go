package main

import (
	"strconv"

	"github.com/faiface/pixel"
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
}
