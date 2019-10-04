package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	Level = level{}
)

type leveler interface {
	update(dt float64, win *pixelgl.Window) leveler
	draw(target pixel.Target)
}

// level is probably going to be the only level
type level struct {
}

func (l *level) update(dt float64, win *pixelgl.Window) leveler {
	if win.JustPressed(pixelgl.KeyU) {
		return &UpgradeScreen
	}

	return &Level
}

func (l *level) draw(target pixel.Target) {

}
