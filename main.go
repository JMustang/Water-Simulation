package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
	title        = "💦 Water Simulation with Raylib-go"
)

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
