package main

import (
	"fmt"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	up    = pixel.V(0, 1)
	down  = pixel.V(0, -1)
	left  = pixel.V(-1, 0)
	right = pixel.V(1, 0)
)

var (
	Enemies []*enemy

	enemy1Sprites []*pixel.Sprite
)

func updateEnemies(dt float64, win *pixelgl.Window) leveler {
	for _, e := range Enemies {
		_ = e.update(dt, win)
	}
	return currentLvl
}

func drawEnemeies(target pixel.Target) {
	for _, e := range Enemies {
		e.draw(target)
	}
}

type enemy struct {
	id              int
	static          bool
	pos             pixel.Vec
	size            pixel.Vec
	searchRange     float64
	attackRange     float64
	sprites         []*pixel.Sprite
	speed           float64
	attackSpeed     float64
	idleThreshold   float64
	requiredUpgrade *upgrade
	lastDir         pixel.Vec
}

// TODO take params not hard coded test values
func NewEnemy() *enemy {
	e := enemy{
		id:              uniqueID(),
		static:          false,
		pos:             camPos.Add(pixel.V(150, 150)),
		size:            pixel.V(16, 16),
		searchRange:     128,
		attackRange:     4,
		sprites:         enemy1Sprites,
		speed:           16 * 3,
		attackSpeed:     16 * 4,
		idleThreshold:   .9,
		requiredUpgrade: movementControls,
	}

	fmt.Println(e.pos)

	return &e
}

func (e *enemy) update(dt float64, win *pixelgl.Window) leveler {
	if !e.requiredUpgrade.acquired {
		return currentLvl
	}

	if pixel.C(e.pos, e.attackRange).IntersectRect(Player.collisionBox()) != pixel.ZV {
		// Can attack
		return currentLvl
	}

	if e.static {
		return currentLvl
	}

	if pixel.C(e.pos, e.searchRange).IntersectRect(Player.collisionBox()) != pixel.ZV {
		// Can see player
		e.moveToPlayer(dt)

		return currentLvl
	}

	if rand.Float64() > e.idleThreshold {
		e.randomWalk(false, dt)
	}

	return currentLvl
}

func (e *enemy) draw(target pixel.Target) {
	if !e.requiredUpgrade.acquired {
		return
	}

	// ToDO animate
	e.sprites[0].Draw(target, pixel.IM.Moved(e.pos))
}

func (e enemy) collisionBox() pixel.Rect {
	return pixel.R(
		e.pos.X,
		e.pos.Y,
		e.pos.Add(e.size).X,
		e.pos.Add(e.size).Y,
	)
}

func (e *enemy) moveToPlayer(dt float64) {
	toPlayer := pixel.L(e.pos, Player.collisionBox().Center())
	toPlayerV := e.pos.To(Player.collisionBox().Center()).Unit().Scaled(e.attackSpeed * dt)
	if !lineCollides(toPlayer) && !rectCollides(e.collisionBox().Moved(toPlayerV)) {
		e.pos = e.pos.Add(toPlayerV)
		return
	}

	e.randomWalk(true, dt)
}

func (e *enemy) randomWalk(attacking bool, dt float64) {
	switch i := rand.Float64(); {
	case i < 0.005:
		e.lastDir = up
	case i < .01:
		e.lastDir = down
	case i < .015:
		e.lastDir = left
	case i < .02:
		e.lastDir = right
	}

	s := e.speed * dt
	if attacking {
		s = e.attackSpeed * dt
	}

	dir := e.lastDir.Scaled(s)

	if !rectCollides(e.collisionBox().Moved(dir)) {
		e.pos = e.pos.Add(dir)
	}
}
