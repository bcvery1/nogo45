package main

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
)

var (
	acquiredUpgrades []*upgrade
)

var (
	upgradeBorderColour              = color.RGBA{R: 0x09, G: 0x0d, B: 0x18, A: 0xff}
	upgradeBackingColour             = color.RGBA{R: 0xb0, G: 0xbb, B: 0xdd, A: 0xff}
	upgradeBackingHoverColour        = color.RGBA{R: 0x52, G: 0xe2, B: 0x52, A: 0xff}
	upgradeBackingHoverColourBlocked = color.RGBA{R: 0x83, G: 0x0, B: 0x00, A: 0xff}
)

var (
	movementControls = &upgrade{
		id:   uniqueID(),
		name: "Movement controls",
		desc: "Gives the players character the ability to move",
		cost: 1,
		next: []*upgrade{seeLevel},
		after: func() {
			DialoguePresenter.queue(afterMovement)
		},
	}

	seeLevel = &upgrade{
		id:   uniqueID(),
		name: "Level",
		desc: "Allows the player to see the level which includes coins",
		cost: 4,
		next: []*upgrade{slowEnemies, mediumEnemies, fastEnemies},
	}

	slowEnemies = &upgrade{
		id:   uniqueID(),
		name: "Basic enemies",
		desc: "Add slow moving enemies to the map.\nEnemies drop coins on death",
		cost: 5,
		next: []*upgrade{},
	}

	mediumEnemies = &upgrade{
		id:   uniqueID(),
		name: "Regular enemies",
		desc: "Add enemies to the map which can move a bit faster.\nEnemies drop coins on death",
		cost: 25,
		next: nil,
	}

	fastEnemies = &upgrade{
		id:   uniqueID(),
		name: "Advanced enemies",
		desc: "Add fast moving enemies to the map.\nEnemies drop coins on death",
		cost: 35,
		next: nil,
	}
)

type upgrade struct {
	id       int
	name     string
	desc     string
	cost     int
	acquired bool
	next     []*upgrade
	after    func()
}

func (u upgrade) draw(target pixel.Target, ind, hoveringOn int) {
	backingC := upgradeBackingColour

	dims := upgradeIndToPos(ind)
	if ind == hoveringOn && u.cost <= Player.coins {
		backingC = upgradeBackingHoverColour
	} else if ind == hoveringOn {
		backingC = upgradeBackingHoverColourBlocked
	}

	imd := imdraw.New(nil)

	imd.Push(cam.Unproject(dims.Min), cam.Unproject(dims.Max))
	imd.Color = backingC
	imd.Rectangle(0)

	grown := pixel.R(
		dims.Min.X-2,
		dims.Min.Y-2,
		dims.Max.X+2,
		dims.Max.Y+2,
	)

	imd.Push(cam.Unproject(grown.Min), cam.Unproject(grown.Max))
	imd.Color = upgradeBorderColour
	imd.Rectangle(3)

	imd.Draw(target)

	moveTitle := dims.Center().Add(pixel.V(-140, 120))
	t := text.New(cam.Unproject(moveTitle), atlas)
	t.Color = defaultTextColour
	_, _ = t.WriteString(fmt.Sprintf("%s - %d coins", u.name, u.cost))
	t.Draw(target, pixel.IM)

	moveDesc := dims.Center().Add(pixel.V(-180, 80))
	t = text.New(cam.Unproject(moveDesc), atlas)
	t.Color = defaultTextColour
	_, _ = t.WriteString(u.desc)
	t.Draw(target, pixel.IM)
}

func (u *upgrade) acquire() {
	if u.cost > Player.coins {
		return
	}
	addCoins(-1 * u.cost)

	u.acquired = true
	acquiredUpgrades = append(acquiredUpgrades, u)
	if u.after != nil {
		u.after()
	}
}

func availableUpgrades() []*upgrade {
	var availUps []*upgrade

	for _, u := range acquiredUpgrades {
		for _, n := range u.next {
			if !n.acquired {
				availUps = append(availUps, n)
			}
		}
	}

	if len(availUps) == 0 {
		return []*upgrade{movementControls}
	}

	return dedup(availUps)
}

func dedup(ups []*upgrade) []*upgrade {
	return ups
}
