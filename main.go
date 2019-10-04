package main

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	winBounds                = pixel.R(0, 0, 1024, 720)
	currentLvl       leveler = &Level
	backgroundColour         = colornames.Whitesmoke
	camPos                   = pixel.ZV
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "NoGo45",
		Bounds: winBounds,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	last := time.Now()

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Moved(winBounds.Center().Sub(camPos))
		win.SetMatrix(cam)

		win.Clear(backgroundColour)

		nextLvl := currentLvl.update(dt, win)
		currentLvl.draw(win)

		win.Update()

		currentLvl = nextLvl
	}
}

func main() {
	pixelgl.Run(run)
}
