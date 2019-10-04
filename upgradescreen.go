package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	UpgradeScreen = upgradeScreen{}
)

type upgradeScreen struct{}

func (u *upgradeScreen) update(dt float64, win *pixelgl.Window) leveler {
	panic("implement me")
}

func (u *upgradeScreen) draw(target pixel.Target) {
	panic("implement me")
}
