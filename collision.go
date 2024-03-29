package main

import (
	"github.com/faiface/pixel"
)

var (
	collisionRs []pixel.Rect

	debugOverride bool
)

func rectCollides(r pixel.Rect) bool {
	if debugOverride || jetpack.acquired {
		return false
	}

	for _, col := range collisionRs {
		if col.Intersect(r) != pixel.ZR {
			return true
		}
	}
	return false
}

func pointCollides(p pixel.Vec) bool {
	for _, col := range collisionRs {
		if col.Contains(p) {
			return true
		}
	}
	return false
}

func lineCollides(l pixel.Line) bool {
	for _, col := range collisionRs {
		if col.IntersectLine(l) != pixel.ZV {
			return true
		}
	}
	return false
}
