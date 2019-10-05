package main

import (
	"image/color"
	"path/filepath"

	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	Level = level{}

	tmxMap     *tilepix.Map
	tilemapPic pixel.Picture
)

func init() {
	var err error
	tmxMap, err = tilepix.ReadFile(filepath.Join(binPath, "assets/level1.tmx"))
	if err != nil {
		panic(err)
	}
}

type leveler interface {
	update(dt float64, win *pixelgl.Window) leveler
	draw(target pixel.Target)
}

// level is probably going to be the only level
type level struct {
}

func (l *level) update(dt float64, win *pixelgl.Window) leveler {
	if isPaused() {
		return currentLvl
	}

	if win.JustPressed(pixelgl.KeyU) {
		firstOpen.Do(func() {
			addCoins(1)
		})
		return &UpgradeScreen
	}

	return l
}

func (l *level) draw(target pixel.Target) {
	if !seeLevel.acquired {
		return
	}

	if err := tmxMap.DrawAll(target, color.Transparent, pixel.IM); err != nil {
		panic(err)
	}
}
