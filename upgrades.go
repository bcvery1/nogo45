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

	movementControls upgrade
	seeLevel         upgrade
	slowEnemies      upgrade
	mediumEnemies    upgrade
	fastEnemies      upgrade
	basicAttack      upgrade
)

var (
	upgradeBorderColour              = color.RGBA{R: 0x09, G: 0x0d, B: 0x18, A: 0xff}
	upgradeBackingColour             = color.RGBA{R: 0xb0, G: 0xbb, B: 0xdd, A: 0xff}
	upgradeBackingHoverColour        = color.RGBA{R: 0x52, G: 0xe2, B: 0x52, A: 0xff}
	upgradeBackingHoverColourBlocked = color.RGBA{R: 0x83, G: 0x0, B: 0x00, A: 0xff}
)

func init() {
	movementControls = upgrade{
		id:   uniqueID(),
		name: "Movement controls",
		desc: "Gives the players character the ability to move",
		cost: 1,
		next: []*upgrade{&seeLevel},
		after: func() {
			DialoguePresenter.queue(afterMovement)
		},
	}

	seeLevel = upgrade{
		id:   uniqueID(),
		name: "Level",
		desc: "Allows the player to see the level which includes coins",
		cost: 4,
		next: []*upgrade{&slowEnemies, &mediumEnemies},
		after: func() {
			for _, obj := range tmxMap.GetObjectByName("coin") {
				p, err := obj.GetPoint()
				if err != nil {
					panic(err)
				}

				Coins = append(Coins, &coin{pos: p})
			}
		},
	}

	slowEnemies = upgrade{
		id:   uniqueID(),
		name: "Basic enemies",
		desc: "Add slow moving enemies to the map.\nEnemies drop coins on death",
		cost: 15,
		next: []*upgrade{&basicAttack},
		after: func() {
			for _, e := range tmxMap.GetObjectByName("e11") {
				p, err := e.GetPoint()
				if err != nil {
					panic(err)
				}
				NewEnemy(p, 1, 1)
			}
			for _, e := range tmxMap.GetObjectByName("e21") {
				p, err := e.GetPoint()
				if err != nil {
					panic(err)
				}
				NewEnemy(p, 2, 1)
			}
			for _, e := range tmxMap.GetObjectByName("e31") {
				p, err := e.GetPoint()
				if err != nil {
					panic(err)
				}
				NewEnemy(p, 3, 1)
			}
		},
	}

	mediumEnemies = upgrade{
		id:   uniqueID(),
		name: "Regular enemies",
		desc: "Add enemies to the map which can move a bit faster.\nEnemies drop coins on death",
		cost: 60,
		next: []*upgrade{&fastEnemies},
		after: func() {
			for _, e := range tmxMap.GetObjectByName("e12") {
				p, err := e.GetPoint()
				if err != nil {
					panic(err)
				}
				NewEnemy(p, 1, 2)
			}
			for _, e := range tmxMap.GetObjectByName("e22") {
				p, err := e.GetPoint()
				if err != nil {
					panic(err)
				}
				NewEnemy(p, 2, 2)
			}
			for _, e := range tmxMap.GetObjectByName("e32") {
				p, err := e.GetPoint()
				if err != nil {
					panic(err)
				}
				NewEnemy(p, 3, 2)
			}
		},
	}

	fastEnemies = upgrade{
		id:   uniqueID(),
		name: "Advanced enemies",
		desc: "Add fast moving enemies to the map.\nEnemies drop coins on death",
		cost: 150,
		next: nil,
		after: func() {
			for _, e := range tmxMap.GetObjectByName("e13") {
				p, err := e.GetPoint()
				if err != nil {
					panic(err)
				}
				NewEnemy(p, 1, 3)
			}
			for _, e := range tmxMap.GetObjectByName("e23") {
				p, err := e.GetPoint()
				if err != nil {
					panic(err)
				}
				NewEnemy(p, 2, 3)
			}
			for _, e := range tmxMap.GetObjectByName("e33") {
				p, err := e.GetPoint()
				if err != nil {
					panic(err)
				}
				NewEnemy(p, 3, 3)
			}
		},
	}

	basicAttack = upgrade{
		id:   uniqueID(),
		name: "Basic strike",
		desc: "Allows the player to strike closeby enemies with mouse 1",
		cost: 25,
	}
}

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
	PlaySound(coinPickupSound)

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
			if n != nil && !n.acquired {
				availUps = append(availUps, n)
			}
		}
	}

	if len(availUps) == 0 {
		return []*upgrade{&movementControls}
	}

	return dedup(availUps)
}

func dedup(ups []*upgrade) []*upgrade {
	return ups
}
