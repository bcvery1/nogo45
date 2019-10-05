package main

import (
	"image/color"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

var (
	winBounds          = pixel.R(0, 0, 1024, 720)
	currentLvl leveler = &Level

	backgroundColour  = colornames.Whitesmoke
	defaultTextColour = color.RGBA{R: 0x10, G: 0x00, B: 0x00, A: 0xff}

	camPos = pixel.ZV
	cam    pixel.Matrix

	atlas *text.Atlas

	binPath string
)

func init() {
	logrus.SetLevel(logrus.FatalLevel)

	bin, err := os.Executable()
	if err != nil {
		panic(err)
	}
	binPath = filepath.Dir(bin)
	tilemapPic, err = loadPicture(filepath.Join(binPath, "assets/tilesheet.png"))
	if err != nil {
		panic(err)
	}

	coinPic = pixel.NewSprite(tilemapPic, spritePos(9, 0))

	playerPics = []*pixel.Sprite{
		pixel.NewSprite(tilemapPic, spritePos(0, 0)),
	}

	enemy1Sprites = []*pixel.Sprite{
		pixel.NewSprite(tilemapPic, spritePos(2, 1)),
	}

	rand.Seed(time.Now().UnixNano())
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "NoGo45",
		Bounds: winBounds,
	}

	atlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	objs := tmxMap.GetObjectByName("start")
	startPoint, err := objs[0].GetPoint()
	if err != nil {
		panic(err)
	}
	camPos = startPoint

	// TODO remove this debug line
	Enemies = append(Enemies, NewEnemy())

	last := time.Now()

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam = pixel.IM.Moved(winBounds.Center().Sub(camPos))
		win.SetMatrix(cam)

		win.Clear(backgroundColour)

		nextLvl := currentLvl.update(dt, win)
		currentLvl.draw(win)

		_ = updateEnemies(dt, win)
		drawEnemeies(win)

		_ = HUD.update(dt, win)
		HUD.draw(win)

		_ = DialoguePresenter.update(dt, win)
		DialoguePresenter.draw(win)

		debug1(win)

		win.Update()

		currentLvl = nextLvl
	}
}

func main() {
	pixelgl.Run(run)
}
