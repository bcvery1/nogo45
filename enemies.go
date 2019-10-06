package main

import (
	"image/color"
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
	for i, e := range Enemies {
		e.update(dt, win, i)
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
	health      float64

	attackRange   float64
	attackTimeout time.Duration
	attackFunc    func(e enemy)
	lastAttack    time.Time
	attackSpeed   float64
	attackColour  color.RGBA
	projSpeed     float64
	attackDam     float64

	sprites []*pixel.Sprite
	static  bool

	requiredUpgrade upgrade
}

// 1 - melee
// 2 - ranged
// 3 - tracking
func NewEnemy(pos pixel.Vec, t, lvl int) {
	e := enemy{
		id:  uniqueID(),
		pos: pos,
	}

	switch t {
	case 1:
		e.attackRange = 4
		e.attackFunc = meleeAttack
		e.attackColour = color.RGBA{
			R: 0x20,
			G: 0x00,
			B: 0x02,
			A: 0xaa,
		}
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
		e.attackColour = color.RGBA{
			R: 0x28,
			G: 0x00,
			B: 0x02,
			A: 0xaa,
		}
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
		e.attackColour = color.RGBA{
			R: 0x3a,
			G: 0x00,
			B: 0x03,
			A: 0xaa,
		}
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
		e.projSpeed = 16 * 4
		e.attackDam = 6
		e.requiredUpgrade = slowEnemies
		e.health = 25
	case 2:
		e.speed = 16 * 4
		e.searchRange = 180
		e.attackTimeout = time.Millisecond * 900
		e.attackSpeed = 16 * 5
		e.projSpeed = 16 * 7
		e.attackDam = 12
		e.requiredUpgrade = mediumEnemies
		e.health = 35
	case 3:
		e.speed = 16 * 5
		e.searchRange = 240
		e.attackTimeout = time.Millisecond * 700
		e.attackSpeed = 16 * 6
		e.projSpeed = 16 * 8
		e.attackDam = 20
		e.requiredUpgrade = fastEnemies
		e.health = 50
	}

	Enemies = append(Enemies, &e)
}

func (e *enemy) update(dt float64, win *pixelgl.Window, ind int) {
	if e == nil {
		return
	}

	if !e.requiredUpgrade.acquired {
		return
	}

	if e.health <= 0 {
		// Dead remove
		copy(Enemies[ind:], Enemies[ind+1:])
		Enemies[len(Enemies)-1] = nil
		Enemies = Enemies[:len(Enemies)-1]
		return
	}

	// Skip if off screen
	if !winBounds.Contains(cam.Project(e.pos)) {
		return
	}

	if pixel.C(e.pos, e.attackRange).IntersectRect(Player.collisionBox()) != pixel.ZV {
		// Can attack
		e.attack()
		return
	}

	if pixel.C(e.pos, e.searchRange).IntersectRect(Player.collisionBox()) != pixel.ZV {
		// Can see player
		if e.static {
			e.aimAtPlayer()
		} else {
			e.moveToPlayer(dt)
		}

		return
	}

	if r := rand.Float64(); r > 0.95 {
		e.randomWalk(false, dt, r)
	}

	return
}

func (e *enemy) draw(target pixel.Target) {
	if !e.requiredUpgrade.acquired {
		return
	}

	e.sprites[0].Draw(target, pixel.IM.Moved(e.pos).Rotated(e.pos, e.angle))
}

func (e *enemy) aimAtPlayer() {
	e.angle = e.pos.To(Player.collisionBox().Center()).Angle()
}

func (e *enemy) attack() {
	if e.lastAttack.Add(e.attackTimeout).Before(time.Now()) {
		e.lastAttack = time.Now()
		e.attackFunc(*e)
	}
}

func (e *enemy) hurt(delta float64) {
	e.health -= delta
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

	e.randomWalk(true, dt, rand.Float64())
}

func (e *enemy) randomWalk(attacking bool, dt float64, r float64) {
	switch {
	case r > 0.999:
		e.lastDir = up
		e.angle = math.Pi / 2
	case r > 0.998:
		e.lastDir = down
		e.angle = (math.Pi * 3) / 2
	case r > 0.997:
		e.lastDir = left
		e.angle = math.Pi
	case r > 0.996:
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
	PlaySound(attackSound)
	NewSwipe(e.pos, e.pos.To(Player.pos()), 16, color.RGBA{
		R: 0xaa,
		G: 0x33,
		B: 0x33,
		A: 0xaa,
	})
	Player.hurt(e.attackDam)
}

func rangedAttack(e enemy) {
	PlaySound(projectile1Sound)
	NewSwipe(e.pos, e.pos.To(Player.pos()), 8, color.RGBA{
		R: 0xff,
		G: 0x88,
		B: 0x88,
		A: 0xff,
	})
	NewProjectile(e.pos, e.pos.To(Player.pos()), e.projSpeed, e.attackDam, 6, e.attackColour, true)
}

func trackingAttack(e enemy) {
	PlaySound(projectile2Sound)
	NewSwipe(e.pos, e.pos.To(Player.pos()), 8, color.RGBA{
		R: 0xff,
		G: 0x88,
		B: 0x88,
		A: 0xff,
	})
	NewProjectile(e.pos, e.pos.To(Player.pos()), e.projSpeed, e.attackDam, 6, e.attackColour, false)
}
