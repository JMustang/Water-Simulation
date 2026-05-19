package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
	title        = "💦 Water Simulation with Raylib-go"
	gridWidth    = 160
	gridHeight   = 90
	cellHeight   = screenHeight / gridHeight //Size of each cell visually
	cellWidth    = screenWidth / gridWidth
	spreadFactor = 0.05                           // How quickly water spreads
	dampening    = 0.999                          // Reduces wave energy over time
	waterColor   = rl.NewColor(50, 100, 200, 255) // A nice blue
	bedrockColor = rl.NewColor(80, 80, 80, 255)
)

type WaterCell struct {
	height   float32
	velocity float32
}

// Global grid for simplicity
var waterGrid [gridWidth][gridHeight]WaterCell
var newWaterGrid [gridWidth][gridHeight]WaterCell

func main() {
	rl.InitWindow(screenWidth, screenHeight, title)
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		//Detect window close button or ESC key
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black) // Clear the background to white
		rl.DrawText("Go Water!", 100, 200, 20, rl.LightGray)
		rl.EndDrawing()
	}
	rl.CloseWindow()
}
