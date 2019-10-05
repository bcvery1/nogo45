package main

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
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

	afterMovement = dialogue{
		text: "Now you can move your character\nUse Esc to get out the upgrade menu\nthen try it out with WASD",
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

	// Box
	diaIMD.Clear()

	diaIMD.Color = colornames.Beige

	bottomLeft := cam.Unproject(pixel.V(diaBorder, diaBottom))
	topLeft := cam.Unproject(pixel.V(diaBorder, diaHeight))
	topRight := cam.Unproject(pixel.V(winBounds.Max.X-diaBorder, diaHeight))

	diaIMD.Push(bottomLeft, topRight)
	diaIMD.Rectangle(0)

	diaIMD.Draw(target)

	// Text
	t := text.New(topLeft.Add(pixel.V(diaPadding, -diaPadding)), atlas)
	t.Color = defaultTextColour
	_, _ = fmt.Fprint(t, d.currentRun.text)
	t.Draw(target, pixel.IM)

	t = text.New(bottomLeft.Add(pixel.V(winBounds.Center().X, diaPadding)), atlas)
	t.Color = defaultTextColour
	_, _ = fmt.Fprint(t, "Press space to continue...")
	t.Draw(target, pixel.IM)
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

// queue queues up a dialogue
func (d *dialoguePresent) queue(dia dialogue) {
	d.currentRun = &dia
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
