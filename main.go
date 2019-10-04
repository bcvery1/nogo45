package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	winBounds = pixel.R(0, 0, 1024, 720)
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

	for !win.Closed() {
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
