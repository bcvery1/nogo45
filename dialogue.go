package main

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	DialoguePresenter = dialoguePresent{
		currentRun: &intro1,
		showing:    false,
	}

	diaIMD = imdraw.New(nil)
)

var (
	intro1 = dialogue{
		text:  "Welcome to NoGo45, a game about starting with nothing",
		next:  &intro2,
		initd: time.Now(),
		delay: time.Second * 3,
	}

	intro2 = dialogue{
		text:  "Oh!  You seem to be missing a character",
		next:  &intro3,
		delay: time.Second * 3,
	}

	intro3 = dialogue{
		text:  "...and controls, oops!",
		next:  &intro4,
		delay: time.Second,
	}

	intro4 = dialogue{
		text:  "Don't worry, we can fix this; we just need coins!",
		next:  &intro5,
		delay: time.Second * 2,
	}

	intro5 = dialogue{
		text:  "To get your first coin, just press 'u' to get to the upgrade menu\nExit this dialogue first with 'Space'",
		delay: time.Second,
	}
)

type dialogue struct {
	text  string
	next  *dialogue
	initd time.Time
	delay time.Duration
}

type dialoguePresent struct {
	currentRun *dialogue
	showing    bool
}

func (d *dialoguePresent) update(dt float64, win *pixelgl.Window) leveler {
	if d.showing {
		d.updateShowing(dt, win)
	} else {
		d.updateNotShowing(dt, win)
	}

	return currentLvl
}

func (d *dialoguePresent) draw(target pixel.Target) {
	if !d.showing {
		return
	}

	const diaBorder = 50.0
	const diaBottom = 150
	const diaHeight = 400
	const diaPadding = 40

	bottomLeftV := pixel.ZV.Sub(cam.Unproject(winBounds.Center())).Sub(winBounds.Center())

	// Box
	diaIMD.Clear()

	diaIMD.Color = colornames.Beige

	bottomLeft := pixel.V(diaBorder, diaBottom).Add(bottomLeftV)
	topLeft := pixel.V(diaBorder, diaHeight).Add(bottomLeftV)
	topRight := pixel.V(winBounds.Max.X-diaBorder, diaHeight).Add(bottomLeftV)

	diaIMD.Push(bottomLeft, topRight)
	diaIMD.Rectangle(0)

	diaIMD.Draw(target)

	// Text
	t.Clear()
	_, _ = fmt.Fprint(t, d.currentRun.text)
	t.Draw(target, pixel.IM.Moved(topLeft.Add(pixel.V(diaPadding, -diaPadding))))

	t.Clear()
	_, _ = fmt.Fprint(t, "Press space to continue...")
	t.Draw(target, pixel.IM.Moved(bottomLeft.Add(pixel.V(winBounds.Center().X, diaPadding))))
}

func (d *dialoguePresent) updateShowing(dt float64, win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeySpace) {
		d.showing = false

		next := d.currentRun.next
		if next != nil {
			next.initd = time.Now()
		}
		d.currentRun = next

		resumeGame()

		return
	}
}

func (d *dialoguePresent) updateNotShowing(dt float64, win *pixelgl.Window) {
	if d.currentRun == nil {
		return
	}

	if time.Since(d.currentRun.initd) >= d.currentRun.delay {
		d.showing = true
		pauseGame()
	}
}
