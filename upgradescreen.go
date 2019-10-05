package main

import (
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	upMargin  = 40
	upPadding = 30
)

var (
	UpgradeScreen = upgradeScreen{
		hoveringOn: -1,
	}

	upWidth  = (winBounds.Max.X - upMargin*3) / 2
	upHeight = (winBounds.Max.Y - upMargin*3) / 2

	upBottomLeftR  = pixel.R(upMargin, upMargin, upMargin+upWidth, upMargin+upHeight)
	upBottomRightR = pixel.R(upMargin*2+upWidth, upMargin, upMargin*2+upWidth*2, upMargin+upHeight)
	upTopLeftR     = pixel.R(upMargin, upMargin*2+upHeight, upMargin+upWidth, upMargin*2+upHeight*2)
	upTopRightR    = pixel.R(upMargin*2+upWidth, upMargin*2+upHeight, upMargin*2+upWidth*2, upMargin*2+upHeight*2)

	firstOpen sync.Once
)

type upgradeScreen struct {
	hoveringOn int
	avail      []*upgrade
}

func (u *upgradeScreen) update(dt float64, win *pixelgl.Window) leveler {
	if win.JustPressed(pixelgl.KeyEscape) {
		return &Level
	}

	if u.avail == nil {
		u.avail = availableUpgrades()
	}

	u.hoveringOn = posToUpgradeInd(win.MousePosition())

	if win.JustPressed(pixelgl.MouseButton1) {
		if u.hoveringOn > -1 && u.hoveringOn < len(u.avail) {
			upClicked := u.avail[u.hoveringOn]
			if upClicked.cost <= Player.coins {
				// Can purchase upgrade
				upClicked.acquire()

				// a panel was clicked
				u.avail = availableUpgrades()
			}
		}
	}

	return u
}

func (u *upgradeScreen) draw(target pixel.Target) {
	for i, up := range availableUpgrades() {
		up.draw(target, i, u.hoveringOn)
	}
}

//  0 | 1
// ---+---
//  2 | 3
func posToUpgradeInd(pos pixel.Vec) int {
	switch {
	case upTopLeftR.Contains(pos):
		return 0
	case upTopRightR.Contains(pos):
		return 1
	case upBottomLeftR.Contains(pos):
		return 2
	case upBottomRightR.Contains(pos):
		return 3
	default:
		return -1
	}
}

func upgradeIndToPos(ind int) pixel.Rect {
	switch ind {
	case 0:
		return upTopLeftR
	case 1:
		return upTopRightR
	case 2:
		return upBottomLeftR
	case 3:
		return upBottomRightR
	default:
		return pixel.ZR
	}
}
