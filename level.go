package main

import (
	"image/color"
	"path/filepath"
	"sync"

	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	nonCollionOGName = "objs"
)

var (
	Level = level{}

	tmxMap     *tilepix.Map
	tilemapPic pixel.Picture

	firstUp    sync.Once
	firstDown  sync.Once
	firstLeft  sync.Once
	firstRight sync.Once

	// pixels per second
	speed = 16. * 8
)

func init() {
	var err error
	tmxMap, err = tilepix.ReadFile(filepath.Join(binPath, "assets/level1.tmx"))
	if err != nil {
		panic(err)
	}

	for _, l := range tmxMap.TileLayers {
		l.SetStatic(true)
	}

	if err := tmxMap.GenerateTileObjectLayer(); err != nil {
		panic(err)
	}
	for _, og := range tmxMap.ObjectGroups {
		// only get collision groups
		if og.Name == nonCollionOGName {
			continue
		}

		for _, obj := range og.Objects {
			r, err := obj.GetRect()
			if err != nil {
				panic(err)
			}

			collisionRs = append(collisionRs, r)
		}
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

	if lvl := Player.update(dt, win); lvl == &DeathScreen {
		return lvl
	}

	_ = updateEnemies(dt, win)
	updateProjectiles(dt)

	if !movementControls.acquired {
		return currentLvl
	}

	deltaV := pixel.ZV
	// up
	if win.Pressed(pixelgl.KeyW) {
		firstUp.Do(func() {
			addCoins(1)
		})

		deltaV = deltaV.Add(pixel.V(0, speed*dt))
	}
	// down
	if win.Pressed(pixelgl.KeyS) {
		firstDown.Do(func() {
			addCoins(1)
		})

		deltaV = deltaV.Add(pixel.V(0, -speed*dt))
	}
	// Allow vertical movement separate to horizontal
	if seeLevel.acquired && deltaV != pixel.ZV && !rectCollides(Player.collisionBox().Moved(deltaV)) {
		camPos = camPos.Add(deltaV)
	}
	deltaV = pixel.ZV

	// left
	if win.Pressed(pixelgl.KeyA) {
		firstLeft.Do(func() {
			addCoins(1)
		})

		deltaV = deltaV.Add(pixel.V(-speed*dt, 0))
	}
	// right
	if win.Pressed(pixelgl.KeyD) {
		firstRight.Do(func() {
			addCoins(1)
		})

		deltaV = deltaV.Add(pixel.V(speed*dt, 0))
	}

	if seeLevel.acquired && deltaV != pixel.ZV && !rectCollides(Player.collisionBox().Moved(deltaV)) {
		camPos = camPos.Add(deltaV)
	}

	return l
}

func (l *level) draw(target pixel.Target) {
	if seeLevel.acquired {
		if err := tmxMap.DrawAll(target, color.Transparent, pixel.IM); err != nil {
			panic(err)
		}

		drawCoins(target)
	}

	Player.draw(target)

	drawEnemeies(target)
	drawProjectiles(target)
}
