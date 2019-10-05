package main

import (
	"github.com/faiface/pixel"
)

var (
	collisionRs []pixel.Rect
)

func rectCollides(r pixel.Rect) bool {
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
