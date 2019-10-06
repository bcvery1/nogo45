// +build !debug

package main

import "github.com/faiface/pixel/pixelgl"

func debug1(_ *pixelgl.Window) {
	return
}

func debugsound() {
	SetupAudio()
}
