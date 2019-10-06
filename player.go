package main

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	healRate = 0.5
)

var (
	Player = player{
		health:    100,
		attackDam: 15,
		aoe:       15,
	}

	playerPics []*pixel.Sprite

	playerSize = pixel.V(16, 16)
)

type player struct {
	coins     int
	angle     float64
	health    float64
	attackDam float64
	aoe       float64
}

func (p *player) update(dt float64, win *pixelgl.Window) leveler {
	if p.health < 0 {
		pauseGame()
		return &DeathScreen
	}

	if p.health < 100 {
		p.health += healRate * dt
	}

	if p.health > 100 {
		p.health = 100
	}

	// check for coin collisions
	coinCollision()

	p.angle = winBounds.Center().To(win.MousePosition()).Angle()

	// Attack
	if basicAttack.acquired && win.JustPressed(pixelgl.MouseButton1) {
		p.attack(win)
	}
	if win.JustPressed(pixelgl.MouseButton2) {
		if atomBomb.acquired {
			p.Atom(win)
		} else if rocketLauncher.acquired {
			p.rocket(win)
		} else if gun.acquired {
			p.gun(win)
		}
	}

	return currentLvl
}

func (p *player) draw(target pixel.Target) {
	// TODO animate
	pos := p.pos()
	playerPics[0].Draw(target, pixel.IM.Moved(pos).Rotated(pos, p.angle))
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

func (p *player) hurt(delta float64) {
	p.health -= delta
	PlaySound(hurtSound)
}

func (p player) pos() pixel.Vec {
	return cam.Unproject(winBounds.Center().Sub(playerSize))
}

func (p player) attack(win *pixelgl.Window) {
	PlaySound(attackSound)

	clickedPos := win.MousePosition()
	toClick := winBounds.Center().To(clickedPos).Unit().Scaled(p.aoe * 1.5).Add(p.pos())
	aoe := pixel.C(toClick, p.aoe)

	NewSwipe(toClick, winBounds.Center().To(clickedPos), p.aoe, color.RGBA{
		R: 0x66,
		G: 0xbb,
		B: 0x66,
		A: 0xbb,
	})

	for _, e := range Enemies {
		if aoe.IntersectRect(e.collisionBox()) != pixel.ZV {
			e.hurt(p.attackDam)
		}
	}
}

func (p player) gun(win *pixelgl.Window) {
	dir := winBounds.Center().To(win.MousePosition())
	NewProjectile(p.pos().Add(dir.Unit().Scaled(1.1*16)), winBounds.Center().To(win.MousePosition()), 16*8, 25, 4, colornames.Black, true)
}

func (p player) rocket(win *pixelgl.Window) {
	dir := winBounds.Center().To(win.MousePosition())
	NewProjectile(p.pos().Add(dir.Unit().Scaled(1.1*16)), winBounds.Center().To(win.MousePosition()), 16*6, 50, 8, colornames.Black, true)
}

func (p player) Atom(win *pixelgl.Window) {
	fmt.Println("boom")
}

func addCoins(delta int) {
	Player.coins += delta
}
