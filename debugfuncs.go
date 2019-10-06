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

	if win.JustPressed(pixelgl.Key1) {
		fmt.Println("Speed to fast")
		speed = 16 * 20
	}
	if win.JustPressed(pixelgl.Key2) {
		fmt.Println("Speed to normal")
		speed = 16 * 6
	}

	if win.JustPressed(pixelgl.Key3) {
		fmt.Println("collisions off")
		debugOverride = true
	}
	if win.JustPressed(pixelgl.Key4) {
		fmt.Println("collisions on")
		debugOverride = false
	}

	if win.JustPressed(pixelgl.Key0) {
		fmt.Println("killing")
		Player.health = -1
	}
}
