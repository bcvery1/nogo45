package main

import "github.com/faiface/pixel"

var (
	Coins []*coin
)

type coin struct {
	pos pixel.Vec
}

func drawCoins(target pixel.Target) {
	if isPaused() {
		return
	}

	for _, c := range Coins {
		coinPic.Draw(target, pixel.IM.Moved(c.pos))
	}
}

func coinCollision() {
	for i, c := range Coins {
		if Player.collisionBox().Contains(c.pos.Add(pixel.V(8, 8))) {
			// Delete coin
			copy(Coins[i:], Coins[i+1:])
			Coins[len(Coins)-1] = nil
			Coins = Coins[:len(Coins)-1]

			addCoins(5)

			return
		}
	}
}
