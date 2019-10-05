package main

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	Player = player{}

	playerPics []*pixel.Sprite

	playerSize = pixel.V(16, 16)
)

type player struct {
	coins int
	angle float64
}

func (p *player) update(dt float64, win *pixelgl.Window) leveler {
	p.angle = winBounds.Center().To(win.MousePosition()).Angle() - math.Pi/2
	return currentLvl
}

func (p *player) draw(target pixel.Target) {
	// TODO animate
	unproj := cam.Unproject(winBounds.Center().Sub(playerSize))
	playerPics[0].Draw(target, pixel.IM.Moved(unproj).Rotated(unproj, p.angle))
}

func (p player) collisionBox() pixel.Rect {
	centre := cam.Unproject(winBounds.Center().Sub(playerSize))

	return pixel.R(
		centre.X,
		centre.Y,
		centre.X+playerSize.X,
		centre.Y+playerSize.Y,
	)
}

func addCoins(delta int) {
	Player.coins += delta
}
