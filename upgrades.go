package main

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

var (
	acquiredUpgrades []upgrade
	allUpgrades      = []upgrade{
		movementControls,
		seeLevel,
		slowEnemies,
		mediumEnemies,
		fastEnemies,
	}
)

var (
	upgradeBorderColour       = color.RGBA{R: 0x09, G: 0x0d, B: 0x18, A: 0xff}
	upgradeBackingColour      = color.RGBA{R: 0xb0, G: 0xbb, B: 0xdd, A: 0xff}
	upgradeBackingHoverColour = color.RGBA{R: 0x52, G: 0xe2, B: 0x52, A: 0xff}
)

var (
	movementControls = upgrade{
		id:   uniqueID(),
		name: "Movement controls",
		desc: "Gives the players character the ability to move",
		cost: 1,
		next: []*upgrade{&seeLevel},
	}

	seeLevel = upgrade{
		id:   uniqueID(),
		name: "Level",
		desc: "Allows the player to see the level which includes coins",
		cost: 4,
		next: []*upgrade{&slowEnemies, &mediumEnemies, &fastEnemies},
	}

	slowEnemies = upgrade{
		id:   uniqueID(),
		name: "Basic enemies",
		desc: "Add slow moving enemies to the map.  Enemies drop coins on death",
		cost: 5,
		next: []*upgrade{},
	}

	mediumEnemies = upgrade{
		id:   uniqueID(),
		name: "Regular enemies",
		desc: "Add enemies to the map which can move a bit faster.  Enemies drop coins on death",
		cost: 25,
		next: nil,
	}

	fastEnemies = upgrade{
		id:   uniqueID(),
		name: "Advanced enemies",
		desc: "Add fast moving enemies to the map.  Enemies drop coins on death",
		cost: 35,
		next: nil,
	}
)

type upgrade struct {
	id   int
	name string
	desc string
	cost int
	next []*upgrade
}

func (u upgrade) draw(imd *imdraw.IMDraw, ind, hoveringOn int) {
	backingC := upgradeBackingColour

	dims := upgradeIndToPos(ind)
	if ind == hoveringOn {
		backingC = upgradeBackingHoverColour
	}

	imd.Push(dims.Min, dims.Max)
	imd.Color = backingC
	imd.Rectangle(0)

	grown := pixel.R(
		dims.Min.X-2,
		dims.Min.Y-2,
		dims.Max.X+2,
		dims.Max.Y+2,
	)

	imd.Push(grown.Min, grown.Max)
	imd.Color = upgradeBorderColour
	imd.Rectangle(3)
}

func availableUpgrades() []upgrade {
	var availUps []upgrade

	for _, u := range acquiredUpgrades {
		for _, n := range u.next {
			availUps = append(availUps, *n)
		}
	}

	if len(availUps) == 0 {
		return []upgrade{movementControls}
	}

	return dedup(availUps)
}

func dedup(ups []upgrade) []upgrade {
	return ups
}
