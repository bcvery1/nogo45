package main

import (
	"image/color"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

var (
	DeathScreen = deathScreen{}
)

type deathScreen struct{}

func (d deathScreen) update(_ float64, win *pixelgl.Window) leveler {
	if win.JustPressed(pixelgl.KeySpace) {
		win.SetClosed(true)
	}

	win.Clear(color.RGBA{
		R: 0x00,
		G: 0x00,
		B: 0x0d,
		A: 0xff,
	})

	return &d
}

func (d deathScreen) draw(target pixel.Target) {
	pos := pixel.V(380, 540)
	pos = cam.Unproject(pos)
	t := text.New(pos, atlas)
	t.Color = color.RGBA{
		R: 0x93,
		G: 0x00,
		B: 0x08,
		A: 0xff,
	}
	_, _ = t.WriteString("You Died")
	t.Draw(target, pixel.IM.Scaled(t.Bounds().Center(), 6))

	pos = pixel.V(330, 280)
	pos = cam.Unproject(pos)
	t = text.New(pos, atlas)
	_, _ = t.WriteString("Thanks for playing, press space to exit")
	t.Draw(target, pixel.IM.Scaled(t.Bounds().Center(), 3))
}

var (
	Nuke        *nuke
	nukeIMD     = imdraw.New(nil)
	explodeOnce sync.Once
)

type nuke struct {
	pos     pixel.Vec
	size    float64
	started time.Time
}

// atom launches the atom bomb, effectively ending the game
func atom(win *pixelgl.Window) {
	if Nuke != nil {
		return
	}

	n := nuke{
		pos:     cam.Unproject(win.MousePosition()),
		started: time.Now(),
		size:    1,
	}

	Nuke = &n
	PlaySound(rumbleSound)

	nukeIMD.Color = color.RGBA{
		R: 0xff,
		G: 0xf5,
		B: 0xae,
		A: 0xff,
	}
}

func updateAtom(dt float64) leveler {
	if Nuke == nil {
		return currentLvl
	}

	// rumble is 4 seconds, wait
	if Nuke.started.Add(time.Second * 4).After(time.Now()) {
		return currentLvl
	}

	explodeOnce.Do(func() {
		PlaySound(explosionSound)
	})

	Nuke.size += Nuke.size * dt * 10

	if Nuke.started.Add(time.Second * 10).Before(time.Now()) {
		return &DeathScreen
	}

	return currentLvl
}

func drawAtom(win *pixelgl.Window) {
	if Nuke == nil {
		return
	}

	nukeIMD.Clear()
	nukeIMD.Push(Nuke.pos)
	nukeIMD.Circle(Nuke.size, 0)
	nukeIMD.Draw(win)
}
