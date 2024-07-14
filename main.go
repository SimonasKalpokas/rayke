package main

// #cgo CFLAGS: -I raylib-5.0_linux_amd64/include
// #cgo LDFLAGS: -L raylib-5.0_linux_amd64/lib -l:libraylib.a -lm
// #include <stdlib.h>
// #include "raylib.h"
import "C"
import "unsafe"

func main() {
	title := C.CString("Title")
	defer C.free(unsafe.Pointer(title))
	C.InitWindow(800, 450, title)

	for !C.WindowShouldClose() {
		C.BeginDrawing()
		C.ClearBackground(C.RAYWHITE)
		C.DrawText(C.CString("Congrats!"), 190, 200, 20, C.LIGHTGRAY)
		C.EndDrawing()
	}
	C.CloseWindow()
}
