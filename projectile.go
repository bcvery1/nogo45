package main

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

var (
	Projectiles []*projectile

	projIMD = imdraw.New(nil)
)

func updateProjectiles(dt float64) {
	for i, p := range Projectiles {
		if p != nil {
			del := p.update(dt)

			if del {
				copy(Projectiles[i:], Projectiles[i+1:])
				Projectiles[len(Projectiles)-1] = nil
				Projectiles = Projectiles[:len(Projectiles)-1]
			}
		}
	}
}

func drawProjectiles(target pixel.Target) {
	projIMD.Clear()

	for _, p := range Projectiles {
		if p != nil {
			p.draw(projIMD)
		}
	}

	projIMD.Draw(target)
}

func NewProjectile(pos, dir pixel.Vec, speed, dam, diameter float64, colour color.RGBA, colide bool) {
	p := projectile{
		pos:        pos,
		dir:        dir.Unit(),
		speed:      speed,
		dam:        dam,
		diameter:   diameter,
		colour:     colour,
		canCollide: colide,
	}

	Projectiles = append(Projectiles, &p)
}

type projectile struct {
	pos        pixel.Vec
	dir        pixel.Vec
	speed      float64
	dam        float64
	diameter   float64
	colour     color.RGBA
	canCollide bool
}

// return if to remove
func (p *projectile) update(dt float64) bool {
	if !winBounds.Contains(cam.Project(p.pos)) {
		// not on screen
		return true
	}

	if Player.collisionBox().Contains(p.pos) {
		Player.hurt(p.dam)
		return true
	}

	if p.canCollide {
		if pointCollides(p.pos) {
			// hit an obsticle
			return true
		}
	}

	// Hit enemies
	for _, e := range Enemies {
		if e == nil {
			continue
		}

		if pixel.C(p.pos, p.diameter/2).IntersectRect(e.collisionBox()) != pixel.ZV {
			e.hurt(p.dam)
			return true
		}
	}

	p.pos = p.pos.Add(p.dir.Scaled(p.speed * dt))

	return false
}

func (p *projectile) draw(imd *imdraw.IMDraw) {
	imd.Color = p.colour
	imd.Push(p.pos)
	imd.Circle(p.diameter/2, 0)
}
