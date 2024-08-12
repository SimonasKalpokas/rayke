package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Direction int

const (
	Left Direction = iota
	Up
	Right
	Down
)

type Coords struct {
	X int32
	Y int32
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
	columnCount int32
	rowCount    int32
	snake       []Coords
	apple       Coords
	direction   Direction
	frameTime   float32
	currentTime float32
	pause       bool

	// rendering
	boxSize int32
}

func New(columnCount, rowCount, boxSize int32, snakeLength int) World {
	frameTime := float32(1) / 10
	world := World{
		columnCount: columnCount,
		rowCount:    rowCount,
		direction:   Right,
		frameTime:   frameTime,
		currentTime: frameTime,
		pause:       false,

		boxSize: boxSize,
	}

	snakeCoords := Coords{X: (columnCount - 1) / 2, Y: (rowCount - 1) / 2}
	for i := 0; i < snakeLength; i++ {
		world.snake = append(world.snake, snakeCoords)
	}

	world.PlaceNewApple()

	return world
}

func (world *World) PlaceNewApple() {
	var x, y int32

	regenerate := true
	for regenerate {
		regenerate = false
		x = rand.Int31n(world.columnCount)
		y = rand.Int31n(world.rowCount)
		for _, snakeCoord := range world.snake {
			if x == snakeCoord.X && y == snakeCoord.Y {
				regenerate = true
				break
			}
		}
	}

	world.apple.X = x
	world.apple.Y = y
}

func (world *World) ScreenHeight() int32 {
	return world.rowCount * world.boxSize
}

func (world *World) ScreenWidth() int32 {
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

	world.snake = append(world.snake, newHead)
	if newHead.X == world.apple.X && newHead.Y == world.apple.Y {
		world.PlaceNewApple()
	} else {
		world.snake = world.snake[1:]
	}
}

func (world *World) Draw() {
	// Draw grid
	for row := int32(0); row < world.rowCount; row++ {
		for column := int32(0); column < world.columnCount; column++ {
			rl.DrawRectangleLines(
				column*world.boxSize,
				row*world.boxSize,
				world.boxSize,
				world.boxSize,
				rl.Gray)
		}
	}

	// Draw snake
	for _, snakeCoords := range world.snake {
		rl.DrawRectangle(
			snakeCoords.X*world.boxSize,
			snakeCoords.Y*world.boxSize,
			world.boxSize,
			world.boxSize,
			rl.Black)
	}

	// Draw apple
	rl.DrawRectangle(
		world.apple.X*world.boxSize,
		world.apple.Y*world.boxSize,
		world.boxSize,
		world.boxSize,
		rl.Red)
}

func main() {
	title := "rayke"

	world := New(40, 20, 20, 5)
	rl.InitWindow(
		world.ScreenWidth(),
		world.ScreenHeight(),
		title)

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeyQ) {
			break
		}
		if rl.IsKeyPressed(rl.KeyP) {
			world.pause = !world.pause
		}
		if rl.IsKeyDown(rl.KeyL) && world.direction != Left {
			world.direction = Right
		}
		if rl.IsKeyDown(rl.KeyH) && world.direction != Right {
			world.direction = Left
		}
		if rl.IsKeyDown(rl.KeyJ) && world.direction != Up {
			world.direction = Down
		}
		if rl.IsKeyDown(rl.KeyK) && world.direction != Down {
			world.direction = Up
		}

		if !world.pause {
			dt := rl.GetFrameTime()
			world.Update(float32(dt))
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.LightGray)
		world.Draw()
		rl.EndDrawing()
	}
	rl.CloseWindow()
}
