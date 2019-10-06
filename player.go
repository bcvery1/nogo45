package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	healRate = 0.5
)

var (
	Player = player{
		health:    100,
		attackDam: 15,
	}

	playerPics []*pixel.Sprite

	playerSize = pixel.V(16, 16)
)

type player struct {
	coins     int
	angle     float64
	health    float64
	attackDam float64
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

	const aoeSize = 15

	clickedPos := win.MousePosition()
	toClick := winBounds.Center().To(clickedPos).Unit().Scaled(aoeSize * 1.5).Add(p.pos())
	aoe := pixel.C(toClick, aoeSize)

	for _, e := range Enemies {
		if aoe.IntersectRect(e.collisionBox()) != pixel.ZV {
			e.hurt(p.attackDam)
		}
	}
}

func addCoins(delta int) {
	Player.coins += delta
}
