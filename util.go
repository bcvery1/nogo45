package main

import "sync/atomic"

var (
	uID int64
)

func uniqueID() int {
	nextID := atomic.AddInt64(&uID, 1)
	return int(nextID)
}
