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

func updateWater() {
	// First, calculate the new state for each cell
	for x := 0; x < gridWidth; x++ {
		for y := 0; y < gridHeight; y++ {
			currentCell := waterGrid[x][y] // Skip if it's considered "solid ground" or empty
			if currentCell.height >= float32(cellHeight)*1.5 || currentCell.height <= 0 {
				newWaterGrid[x][y] = currentCell
				continue
			} // Apply gravity (simple downward pull)
			currentCell.velocity -= 0.05 // A small downward acceleration			// Pressure and flow with neighbors
			// We need to check bounds for neighbors, otherwise our program will crash!
			if x > 0 { // Left neighbor
				diff := currentCell.height - waterGrid[x-1][y].height
				currentCell.velocity -= diff * spreadFactor
				waterGrid[x-1][y].velocity += diff * spreadFactor // Affect neighbor's velocity
			}
			if x < gridWidth-1 { // Right neighbor
				diff := currentCell.height - waterGrid[x+1][y].height
				currentCell.velocity -= diff * spreadFactor
				waterGrid[x+1][y].velocity += diff * spreadFactor
			}
			if y > 0 { // Top neighbor (less common for basic water, but for splash effects)
				diff := currentCell.height - waterGrid[x][y-1].height
				currentCell.velocity -= diff * spreadFactor * 0.5 // Weaker vertical spread
				waterGrid[x][y-1].velocity += diff * spreadFactor * 0.5
			}
			if y < gridHeight-1 { // Bottom neighbor (gravity already handles this, but for upward force)
				// You know, we won't explicitly add bottom neighbor pressure here, gravity really dominates that.
			} // Update height based on velocity
			currentCell.height += currentCell.velocity // Apply damping
			currentCell.velocity *= dampening          // Clamp height to prevent extreme values (e.g., negative water).
			// We don't want our water disappearing or exploding, right?
			if currentCell.height < 0 {
				currentCell.height = 0
			}
			if currentCell.height > float32(cellHeight)*2 { // Don't let it explode
				currentCell.height = float32(cellHeight) * 2
			}
			newWaterGrid[x][y] = currentCell
		}
	} // Copy the new state back to the main grid. This is crucial for the next frame!
	for x := 0; x < gridWidth; x++ {
		for y := 0; y < gridHeight; y++ {
			// A little sanity check: ensure "solid" ground actually stays solid.
			if waterGrid[x][y].height >= float32(cellHeight)*1.5 { // A threshold to indicate solid ground
				waterGrid[x][y].height = float32(cellHeight) * 2
				waterGrid[x][y].velocity = 0
			} else {
				waterGrid[x][y] = newWaterGrid[x][y]
			}
		}
	}
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, title)
	rl.SetTargetFPS(60)

	for x := 0; x < gridWidth; x++ {
		for y := 0; y < gridHeight; y++ {
			if y > gridHeight/2 { // Start with water in the bottom half
				waterGrid[x][y].height = float32(cellHeight)
			} else {
				waterGrid[x][y].height = 0
			}
			waterGrid[x][y].velocity = 0
		}
	}
	// Let's add some ground to the very bottom
	for x := 0; x < gridWidth; x++ {
		for y := gridHeight - 5; y < gridHeight; y++ { // Bottom 5 rows are solid
			waterGrid[x][y].height = float32(cellHeight) * 2 // Give it extra "height" to signify solid
		}
	}

	for !rl.WindowShouldClose() {
		//Detect window close button or ESC key
		updateWater()
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black) // Clear the background to white
		rl.DrawText("Go Water!", 100, 200, 20, rl.LightGray)
		rl.EndDrawing()
	}
	rl.CloseWindow()
}
