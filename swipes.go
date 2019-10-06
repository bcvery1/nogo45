package main

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

const (
	fadeRate float64 = 255 * 8
	moveRate float64 = 16 * 2
)

var (
	Swipes   []*swipe
	swipeIMD = imdraw.New(nil)
)

type swipe struct {
	pos  pixel.Vec
	dir  pixel.Vec
	size float64
	c    color.RGBA
}

func NewSwipe(pos, dir pixel.Vec, size float64, c color.RGBA) {
	s := swipe{
		pos:  pos,
		dir:  dir.Unit(),
		size: size,
		c:    c,
	}

	Swipes = append(Swipes, &s)
}

func updateSwipes(dt float64) {
	for i, s := range Swipes {
		if s.update(dt) {
			copy(Swipes[i:], Swipes[i+1:])
			Swipes[len(Swipes)-1] = nil
			Swipes = Swipes[:len(Swipes)-1]
		}
	}
}

func drawSwipes(target pixel.Target) {
	swipeIMD.Clear()

	for _, s := range Swipes {
		s.draw(swipeIMD)
	}

	swipeIMD.Draw(target)
}

func (s *swipe) update(dt float64) bool {
	if s == nil {
		return false
	}

	nextFade := uint8(dt * fadeRate)
	if s.c.A <= nextFade {
		return true
	}

	if s.c.R < nextFade {
		s.c.R = 0x00
	}
	if s.c.G < nextFade {
		s.c.G = 0x00
	}
	if s.c.B < nextFade {
		s.c.B = 0x00
	}
	s.c.A -= nextFade

	s.pos = s.pos.Add(s.dir.Scaled(dt * moveRate))

	return false
}

func (s *swipe) draw(imd *imdraw.IMDraw) {
	if s == nil {
		return
	}

	imd.Color = s.c
	imd.Push(s.pos)
	imd.Circle(s.size, 0)
}
