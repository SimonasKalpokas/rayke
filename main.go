package main

// #cgo CFLAGS: -I raylib-5.0_linux_amd64/include
// #cgo LDFLAGS: -L raylib-5.0_linux_amd64/lib -l:libraylib.a -lm
// #include <stdlib.h>
// #include "raylib.h"
import "C"
import (
	"unsafe"
)

type Direction int

const (
	Left Direction = iota
	Up
	Right
	Down
)

type Coords struct {
	X int
	Y int
}

func (c Coords) MoveInDirection(d Direction) Coords {
	switch d {
	case Left:
		return Coords{X: c.X - 1, Y: c.Y}
	case Up:
		return Coords{X: c.X, Y: c.Y - 1}
	case Right:
		return Coords{X: c.X + 1, Y: c.Y}
	case Down:
		return Coords{X: c.X, Y: c.Y + 1}
	default:
		panic("Impossible direction")
	}
}

type World struct {
	// logic
	columnCount int
	rowCount    int
	snake       []Coords
	direction   Direction
	frameTime   float32
	currentTime float32

	// rendering
	boxSize int
}

func New(columnCount int, rowCount int, boxSize int, snakeLength int) World {
	frameTime := float32(1) / 10
	world := World{
		columnCount: columnCount,
		rowCount:    rowCount,
		direction:   Right,
		frameTime:   frameTime,
		currentTime: frameTime,

		boxSize: boxSize,
	}

	snakeCoords := Coords{X: (columnCount - 1) / 2, Y: (rowCount - 1) / 2}
	for i := 0; i < snakeLength; i++ {
		world.snake = append(world.snake, snakeCoords)
	}

	return world
}

func (world *World) ScreenHeight() int {
	return world.rowCount * world.boxSize
}

func (world *World) ScreenWidth() int {
	return world.columnCount * world.boxSize
}

func (world *World) Update(dt float32) {
	world.currentTime = world.currentTime - dt
	if world.currentTime > 0 {
		return
	}

	world.currentTime = world.frameTime
	head := world.snake[len(world.snake)-1]
	newHead := head.MoveInDirection(world.direction)
	if newHead.X < 0 {
		newHead.X = world.columnCount - 1
	}
	if newHead.Y < 0 {
		newHead.Y = world.rowCount - 1
	}
	if newHead.X >= world.columnCount {
		newHead.X = 0
	}
	if newHead.Y >= world.rowCount {
		newHead.Y = 0
	}
	world.snake = append(world.snake, newHead)[1:]
}

func (world *World) Draw() {
	for row := 0; row < world.rowCount; row++ {
		for column := 0; column < world.columnCount; column++ {
			C.DrawRectangleLines(
				C.int(column*world.boxSize),
				C.int(row*world.boxSize),
				C.int(world.boxSize),
				C.int(world.boxSize),
				C.GRAY)
		}
	}

	for _, snakeCoords := range world.snake {
		C.DrawRectangle(
			C.int(snakeCoords.X*world.boxSize),
			C.int(snakeCoords.Y*world.boxSize),
			C.int(world.boxSize),
			C.int(world.boxSize),
			C.BLACK)
	}
}

func main() {
	title := C.CString("rayke")
	defer C.free(unsafe.Pointer(title))

	world := New(40, 20, 20, 5)
	C.InitWindow(
		C.int(world.ScreenWidth()),
		C.int(world.ScreenHeight()),
		title)

	for !C.WindowShouldClose() {
		dt := C.GetFrameTime()
		world.Update(float32(dt))

		if C.IsKeyDown(C.KEY_RIGHT) && world.direction != Left {
			world.direction = Right
		}
		if C.IsKeyDown(C.KEY_LEFT) && world.direction != Right {
			world.direction = Left
		}
		if C.IsKeyDown(C.KEY_DOWN) && world.direction != Up {
			world.direction = Down
		}
		if C.IsKeyDown(C.KEY_UP) && world.direction != Down {
			world.direction = Up
		}
		C.BeginDrawing()
		C.ClearBackground(C.LIGHTGRAY)
		world.Draw()
		C.EndDrawing()
	}
	C.CloseWindow()
}
