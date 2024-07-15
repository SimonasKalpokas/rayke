package main

// #cgo CFLAGS: -I raylib-5.0_linux_amd64/include
// #cgo LDFLAGS: -L raylib-5.0_linux_amd64/lib -l:libraylib.a -lm
// #include <stdlib.h>
// #include "raylib.h"
import "C"
import "unsafe"

type World struct {
	// logic
	columnCount int
	rowCount    int

	// rendering
	boxSize int
}

func (world *World) ScreenHeight() int {
	return world.rowCount * world.boxSize
}

func (world *World) ScreenWidth() int {
	return world.columnCount * world.boxSize
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
}

func main() {
	title := C.CString("Title")
	defer C.free(unsafe.Pointer(title))

	world := World{columnCount: 40, rowCount: 20, boxSize: 20}
	C.InitWindow(C.int(world.ScreenWidth()), C.int(world.ScreenHeight()), title)

	for !C.WindowShouldClose() {
		C.BeginDrawing()
		C.ClearBackground(C.LIGHTGRAY)
		world.Draw()
		C.EndDrawing()
	}
	C.CloseWindow()
}
