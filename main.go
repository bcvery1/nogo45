package main

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

var (
	winBounds                = pixel.R(0, 0, 1024, 720)
	currentLvl       leveler = &Level
	backgroundColour         = colornames.Whitesmoke
	camPos                   = pixel.ZV
	cam              pixel.Matrix
	atlas            *text.Atlas
	t                *text.Text
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "NoGo45",
		Bounds: winBounds,
	}

	atlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)
	t = text.New(pixel.ZV, atlas)
	t.Color = colornames.Navy

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	last := time.Now()

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam = pixel.IM.Moved(winBounds.Center().Sub(camPos))
		win.SetMatrix(cam)

		win.Clear(backgroundColour)

		nextLvl := currentLvl.update(dt, win)
		currentLvl.draw(win)

		_ = DialoguePresenter.update(dt, win)
		DialoguePresenter.draw(win)

		// TODO remove this debug
		if win.JustPressed(pixelgl.MouseButton1) {
			mPos := win.MousePosition()
			fmt.Printf("Window (%.0f, %.0f) Projected (%.0f, %.0f)\n", mPos.X, mPos.Y, cam.Unproject(mPos).X, cam.Unproject(mPos).Y)
		}

		win.Update()

		currentLvl = nextLvl
	}
}

func main() {
	pixelgl.Run(run)
}