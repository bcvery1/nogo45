package main

import (
	"fmt"
	"image"
	"os"
	"sync/atomic"

	_ "image/png"

	"github.com/faiface/pixel"
)

var (
	uID int64
)

func uniqueID() int {
	nextID := atomic.AddInt64(&uID, 1)
	return int(nextID)
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

// i and j zero indexed, from bottom left
func spritePos(i, j int) pixel.Rect {
	iF := float64(i)
	jF := float64(j)
	r := pixel.R(
		iF*16,
		jF*16.,
		(iF+1)*16,
		(jF+1)*16,
	)
	fmt.Println(r)
	return r
}
