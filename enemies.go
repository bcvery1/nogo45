package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

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

	enemy11Sprites []*pixel.Sprite
	enemy12Sprites []*pixel.Sprite
	enemy13Sprites []*pixel.Sprite

	enemy21Sprites []*pixel.Sprite
	enemy22Sprites []*pixel.Sprite
	enemy23Sprites []*pixel.Sprite

	enemy31Sprites []*pixel.Sprite
	enemy32Sprites []*pixel.Sprite
	enemy33Sprites []*pixel.Sprite
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
	attackTimeout   time.Duration
	attackFunc      func(e enemy)
	lastAttack      time.Time
	sprites         []*pixel.Sprite
	speed           float64
	attackSpeed     float64
	idleThreshold   float64
	requiredUpgrade *upgrade
	lastDir         pixel.Vec
	angle           float64
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
		sprites:         enemy11Sprites,
		speed:           16 * 3,
		attackSpeed:     16 * 4,
		attackTimeout:   time.Second * 2,
		attackFunc:      meleeAttack,
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
		e.attack()
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
	e.sprites[0].Draw(target, pixel.IM.Moved(e.pos).Rotated(e.pos, e.angle))
}

func (e *enemy) attack() {
	if e.lastAttack.Add(e.attackTimeout).Before(time.Now()) {
		e.lastAttack = time.Now()
		e.attackFunc(*e)
	}
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
		e.angle = toPlayerV.Angle()
		return
	}

	e.randomWalk(true, dt)
}

func (e *enemy) randomWalk(attacking bool, dt float64) {
	switch i := rand.Float64(); {
	case i < 0.005:
		e.lastDir = up
		e.angle = math.Pi / 2
	case i < .01:
		e.lastDir = down
		e.angle = (math.Pi * 3) / 2
	case i < .015:
		e.lastDir = left
		e.angle = math.Pi
	case i < .02:
		e.lastDir = right
		e.angle = 0
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

// **************** Attack functions *************************** \\

func meleeAttack(e enemy) {
	fmt.Println("Melee attack")
}

func rangedAttack(e enemy) {
	fmt.Println("Ranged attack")
}

func trackingAttack(e enemy) {
	fmt.Println("Tracking attack")
}
