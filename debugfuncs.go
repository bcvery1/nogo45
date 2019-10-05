// +build debug

package main

import (
	"fmt"

	"github.com/faiface/pixel/pixelgl"
)

func debug1(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.MouseButton1) {
		mPos := win.MousePosition()
		fmt.Printf("Window (%.0f, %.0f) Projected (%.0f, %.0f)\n", mPos.X, mPos.Y, cam.Unproject(mPos).X, cam.Unproject(mPos).Y)
	}
}
