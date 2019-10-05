package main

import "sync/atomic"

var (
	currentGameState = gamePaused
)

const (
	gamePaused int64 = iota
	gameRunning
)

func isPaused() bool {
	return atomic.LoadInt64(&currentGameState) == gamePaused
}

func pauseGame() {
	atomic.StoreInt64(&currentGameState, gamePaused)
}

func resumeGame() {
	atomic.StoreInt64(&currentGameState, gameRunning)
}
