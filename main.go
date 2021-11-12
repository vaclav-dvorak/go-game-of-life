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

const numCells = 200

var (
	fullscreen    bool    = false
	sizePx        int     = 600
	version, date         = "dev", "now"
	scale         float64 = float64(sizePx) / numCells
)

type Cell struct {
	x, y, decay int
	isAlive     bool
}

var grid []Cell

func getColor() []color.RGBA {
	return []color.RGBA{
		{2, 2, 2, 255},
		{23, 45, 23, 255},
		{44, 87, 44, 255},
		{65, 129, 65, 255},
		{86, 171, 86, 255},
		{107, 213, 107, 255},
		{128, 255, 128, 255},
	}
}

func countGridIndex(x, y int) int {
	return x + y*numCells
}

func isAlive(x, y int) int {
	if x < 0 || x >= numCells || y < 0 || y >= numCells {
		return 0
	}
	if grid[countGridIndex(x, y)].isAlive {
		return 1
	} else {
		return 0
	}
}

func initGrid() []Cell {
	rand.Seed(time.Now().UnixNano())
	start := make([]Cell, numCells*numCells)
	for x := 0; x < numCells; x++ {
		for y := 0; y < numCells; y++ {
			start[countGridIndex(x, y)] = Cell{
				x:       x,
				y:       y,
				isAlive: rand.Float64() > 0.5,
				decay:   0,
			}
			if start[countGridIndex(x, y)].isAlive {
				start[countGridIndex(x, y)].decay = 6
			}
		}
	}
	return start
}

func nextGeneration() []Cell {
	gen := make([]Cell, numCells*numCells)
	for idx, cell := range grid {
		gen[idx] = cell
		around := isAlive(cell.x-1, cell.y-1) +
			isAlive(cell.x, cell.y-1) +
			isAlive(cell.x+1, cell.y-1) +
			isAlive(cell.x-1, cell.y) +
			isAlive(cell.x+1, cell.y) +
			isAlive(cell.x-1, cell.y+1) +
			isAlive(cell.x, cell.y+1) +
			isAlive(cell.x+1, cell.y+1)
		if around == 2 { // Do nothing
			gen[idx].isAlive = cell.isAlive
		} else if around == 3 { // Make alive
			gen[idx].isAlive = true
			gen[idx].decay = 6
		} else { // Make dead
			gen[idx].isAlive = false
			if cell.decay > 0 {
				gen[idx].decay = cell.decay - 1
			} else {
				gen[idx].decay = 0
			}
		}
	}

	return gen
}

func world() *image.RGBA {
	w := image.NewRGBA(image.Rect(0, 0, numCells, numCells))

	for _, cell := range grid {
		w.Set(cell.x, cell.y, getColor()[cell.decay])
	}

	return w
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:       fmt.Sprintf("Game of life %s (%s)", version, date),
		Bounds:      pixel.R(0, 0, float64(numCells)*scale, float64(numCells)*scale),
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

		grid = nextGeneration()
		win.Clear(color.Black)

		w := pixel.PictureDataFromImage(world())

		pixel.NewSprite(w, w.Bounds()).
			Draw(win, pixel.IM.
				Scaled(pixel.ZV, scale).
				Moved(win.Bounds().Center()))

		win.Update()
	}
}

func main() {
	grid = initGrid()
	pixelgl.Run(run)
}
