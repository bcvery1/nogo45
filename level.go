package main

import (
	"image/color"
	"path/filepath"
	"sync"

	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	Level = level{}

	tmxMap     *tilepix.Map
	tilemapPic pixel.Picture

	firstUp    sync.Once
	firstDown  sync.Once
	firstLeft  sync.Once
	firstRight sync.Once
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

	_ = Player.update(dt, win)

	if !movementControls.acquired {
		return currentLvl
	}

	// up
	if win.Pressed(pixelgl.KeyW) {
		firstUp.Do(func() {
			addCoins(1)
		})
	}
	// down
	if win.Pressed(pixelgl.KeyS) {
		firstDown.Do(func() {
			addCoins(1)
		})
	}
	// left
	if win.Pressed(pixelgl.KeyA) {
		firstLeft.Do(func() {
			addCoins(1)
		})
	}
	// right
	if win.Pressed(pixelgl.KeyD) {
		firstRight.Do(func() {
			addCoins(1)
		})
	}

	return l
}

func (l *level) draw(target pixel.Target) {
	Player.draw(target)

	if !seeLevel.acquired {
		return
	}

	if err := tmxMap.DrawAll(target, color.Transparent, pixel.IM); err != nil {
		panic(err)
	}
}
