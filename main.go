package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"log"
	"time"
	"unsafe"
)

const (
	SCREEN_WIDTH  = 1600
	SCREEN_HEIGHT = 1600
	SQUARE_SIZE   = 40.0
)

var (
	grid_x    = int(SCREEN_WIDTH / SQUARE_SIZE)
	grid_y    = int(SCREEN_HEIGHT / SQUARE_SIZE)
	grid      = make([][]bool, grid_x)
	next_grid = make([][]bool, grid_x)
	pause     = false
)

func main() {
	// Game of life init
	for i := range grid {
		grid[i] = make([]bool, grid_y)
		next_grid[i] = make([]bool, grid_y)
	}

	// GLFW init
	err := glfw.Init()
	if err != nil {
		log.Fatalln(err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "Gome Of Life", nil, nil)
	if err != nil {
		log.Fatalln(err)
	}

	window.SetKeyCallback(key_callback)
	window.SetMouseButtonCallback(mouse_callback)

	window.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}

	// Sets up the opengl viewport so that the coordinates correspond to pixels on the screen
	gl.Viewport(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, SCREEN_WIDTH, 0, SCREEN_HEIGHT, 0, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	beginFrame := time.Now()
	var deltaTime time.Time
	timestep := 100

	for !window.ShouldClose() {
		deltaTime = time.Now()
		if int((deltaTime.Sub(beginFrame)))/1000000 > timestep && !pause {
			beginFrame = time.Now()

			// Advance the frame
			for i := 0; i < grid_x; i++ {
				for j := 0; j < grid_y; j++ {
					neighbours :=
						check(i-1, j-1) + check(i, j-1) + check(i+1, j-1) +
							check(i-1, j) + check(i+1, j) +
							check(i-1, j+1) + check(i, j+1) + check(i+1, j+1)

					next_grid[i][j] =
						(grid[i][j] && (neighbours == 2 || neighbours == 3)) ||
							(!grid[i][j] && neighbours == 3)
				}
			}

			for i := 0; i < grid_x; i++ {
				for j := 0; j < grid_y; j++ {
					grid[i][j] = next_grid[i][j]
					next_grid[i][j] = false
				}
			}

			draw(window)

			// Resets opengl context if the context is lost
			if glfw.GetCurrentContext() == nil {
				window.MakeContextCurrent()
			}

		}
		glfw.PollEvents()
	}
}

// Returns 1 if coordinate is in bounds and alive,
// returns 0 if coordinate is out of bounds or dead
func check(x, y int) int {
	var ret_val int
	if x >= 0 && x < grid_x && y >= 0 && y < grid_y && grid[x][y] {
		ret_val = 1
	}
	return ret_val
}

func draw(w *glfw.Window) {
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	for i := 0; i < grid_x; i++ {
		for j := 0; j < grid_y; j++ {
			drawBox(i, j, grid[i][j])
		}
	}

	w.SwapBuffers()
}

var white = [12]float32{
	1, 1, 1,
	1, 1, 1,
	1, 1, 1,
	1, 1, 1}
var grey = [12]float32{
	0.5, 0.5, 0.5,
	0.5, 0.5, 0.5,
	0.5, 0.5, 0.5,
	0.5, 0.5, 0.5}
var black = [12]float32{
	0, 0, 0,
	0, 0, 0,
	0, 0, 0,
	0, 0, 0}

func drawBox(x, y int, alive bool) {
	var adjusted_x float32 = (float32(x) * SQUARE_SIZE)
	var adjusted_y float32 = (float32(y) * SQUARE_SIZE)

	var borderVertices = [8]float32{
		adjusted_x, adjusted_y, adjusted_x,
		adjusted_y + SQUARE_SIZE, adjusted_x + SQUARE_SIZE, adjusted_y + SQUARE_SIZE,
		adjusted_x + SQUARE_SIZE, adjusted_y}
	var fillVertices = [8]float32{
		(adjusted_x) + 2, (adjusted_y) + 2, (adjusted_x) + 2,
		(adjusted_y) + SQUARE_SIZE - 2, (adjusted_x) + SQUARE_SIZE - 2, (adjusted_y) + SQUARE_SIZE - 2,
		(adjusted_x) + SQUARE_SIZE - 2, (adjusted_y) + 2}

	gl.EnableClientState(gl.VERTEX_ARRAY)
	gl.EnableClientState(gl.COLOR_ARRAY)

	gl.VertexPointer(2, gl.FLOAT, 0, unsafe.Pointer(&borderVertices))
	gl.ColorPointer(3, gl.FLOAT, 0, unsafe.Pointer(&grey))
	gl.DrawArrays(gl.POLYGON, 0, 4)

	gl.VertexPointer(2, gl.FLOAT, 0, unsafe.Pointer(&fillVertices))
	if alive {
		gl.ColorPointer(3, gl.FLOAT, 0, unsafe.Pointer(&white))
	} else {
		gl.ColorPointer(3, gl.FLOAT, 0, unsafe.Pointer(&black))
	}
	gl.DrawArrays(gl.POLYGON, 0, 4)

	gl.DisableClientState(gl.COLOR_ARRAY)
	gl.DisableClientState(gl.VERTEX_ARRAY)
}

func key_callback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == 0 {
		return
	}
	if key == 256 || key == 81 { // escape and q
		w.SetShouldClose(true)
	} else if key == 32 { // space
		pause = !pause
	} else if key == 67 { // c
		for i := 0; i < grid_x; i++ {
			for j := 0; j < grid_y; j++ {
				grid[i][j] = false
				next_grid[i][j] = false
			}
		}
		draw(w)
	}
}

func mouse_callback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if action == 0 {
		return
	}
	xpos, ypos := w.GetCursorPos()
	ypos = SCREEN_HEIGHT - ypos
	x := int(xpos / SQUARE_SIZE)
	y := int(ypos / SQUARE_SIZE)

	grid[x][y] = !grid[x][y]

	draw(w)
}
