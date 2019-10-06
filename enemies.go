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
	id int

	pos     pixel.Vec
	speed   float64
	lastDir pixel.Vec
	angle   float64

	searchRange float64

	attackRange   float64
	attackTimeout time.Duration
	attackFunc    func(e enemy)
	lastAttack    time.Time
	attackSpeed   float64
	attackDam     float64

	sprites       []*pixel.Sprite
	static        bool
	idleThreshold float64

	requiredUpgrade upgrade
}

// 1 - melee
// 2 - ranged
// 3 - tracking
func NewEnemy(pos pixel.Vec, t, lvl int) {
	e := enemy{
		id:            uniqueID(),
		pos:           pos,
		idleThreshold: 0.95,
	}

	switch t {
	case 1:
		e.attackRange = 4
		e.attackFunc = meleeAttack
		e.static = false

		e.sprites = enemy11Sprites
		if lvl == 2 {
			e.sprites = enemy12Sprites
		} else if lvl == 3 {
			e.sprites = enemy13Sprites
		}
	case 2:
		e.attackRange = 100
		e.attackFunc = rangedAttack
		e.static = true

		e.sprites = enemy21Sprites
		if lvl == 2 {
			e.sprites = enemy22Sprites
		} else if lvl == 3 {
			e.sprites = enemy23Sprites
		}
	case 3:
		e.attackRange = 80
		e.attackFunc = trackingAttack
		e.static = false

		e.sprites = enemy31Sprites
		if lvl == 2 {
			e.sprites = enemy32Sprites
		} else if lvl == 3 {
			e.sprites = enemy33Sprites
		}
	}

	switch lvl {
	case 1:
		e.speed = 16 * 2
		e.searchRange = 120
		e.attackTimeout = time.Millisecond * 1300
		e.attackSpeed = 16 * 4
		e.requiredUpgrade = slowEnemies
	case 2:
		e.speed = 16 * 4
		e.searchRange = 180
		e.attackTimeout = time.Millisecond * 900
		e.attackSpeed = 16 * 5
		e.requiredUpgrade = mediumEnemies
	case 3:
		e.speed = 16 * 5
		e.searchRange = 240
		e.attackTimeout = time.Millisecond * 700
		e.attackSpeed = 16 * 6
		e.requiredUpgrade = fastEnemies
	}

	Enemies = append(Enemies, &e)
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
		e.pos.Add(pixel.V(16, 16)).X,
		e.pos.Add(pixel.V(16, 16)).Y,
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
