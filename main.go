package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const numCells = 3

var (
	fullscreen    = false
	width         = 320
	height        = 200
	scale         = 3.0
	version, date = "dev", "now"
)

type Cell struct {
	x, y    uint16
	isAlive bool
	decay   uint8
}

var grid []Cell

func countGridIndex(x, y int) uint16 {
	return uint16(x + y*numCells)
}

func initGrid() []Cell {
	rand.Seed(time.Now().UnixNano())
	slice := make([]Cell, numCells*numCells)
	for x := 0; x < numCells; x++ {
		for y := 0; y < numCells; y++ {
			slice[countGridIndex(x, y)] = Cell{
				x:       uint16(x),
				y:       uint16(y),
				isAlive: rand.Float64() > 0.5,
				decay:   0,
			}
		}
	}
	return slice
}

func minimap() *image.RGBA {
	m := image.NewRGBA(image.Rect(0, 0, numCells, numCells))

	// for x, row := range world {
	// 	for y, _ := range row {
	// 		c := getColor(x, y)
	// 		if c.A == 255 {
	// 			c.A = 96
	// 		}
	// 		m.Set(x, y, c)
	// 	}
	// }

	return m
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:       fmt.Sprintf("Game of life %s (%s)", version, date),
		Bounds:      pixel.R(0, 0, float64(width)*scale, float64(height)*scale),
		VSync:       true,
		Undecorated: true,
	}

	if fullscreen {
		cfg.Monitor = pixelgl.PrimaryMonitor()
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			return
		}

		win.Clear(color.Black)

		// m := pixel.PictureDataFromImage(minimap())

		// pixel.NewSprite(m, m.Bounds()).
		// 	Draw(win, pixel.IM)

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
	grid = initGrid()
	fmt.Printf("%+v", grid)
}
