package main

var (
	Player = player{}
)

type player struct {
	coins int
}

func addCoins(delta int) {
	Player.coins += delta
}
