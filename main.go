package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
	title        = "Water Simulation with Raylib-go"
	gridWidth    = 160
	gridHeight   = 90
	cellHeight   = screenHeight / gridHeight
	cellWidth    = screenWidth / gridWidth
	spreadFactor = 0.25  // Wave propagation strength
	dampening    = 0.985 // Energy loss per frame
	gravity      = 0.02
)

var (
	waterColor   = rl.NewColor(50, 100, 200, 255)
	bedrockColor = rl.NewColor(80, 80, 80, 255)
)

type WaterCell struct {
	height   float32
	velocity float32
}

var waterGrid [gridWidth][gridHeight]WaterCell
var newWaterGrid [gridWidth][gridHeight]WaterCell

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func isBedrock(height float32) bool {
	return height >= float32(cellHeight)*1.5
}

func heightAt(x, y int) float32 {
	if x < 0 || x >= gridWidth || y < 0 || y >= gridHeight {
		return 0
	}
	h := waterGrid[x][y].height
	if isBedrock(h) {
		return float32(cellHeight) * 2
	}
	return h
}

func laplacian(x, y int) float32 {
	center := heightAt(x, y)
	return heightAt(x-1, y) + heightAt(x+1, y) + heightAt(x, y-1) + heightAt(x, y+1) - 4*center
}

func updateWater() {
	for x := 0; x < gridWidth; x++ {
		for y := 0; y < gridHeight; y++ {
			cell := waterGrid[x][y]
			if isBedrock(cell.height) {
				newWaterGrid[x][y] = WaterCell{
					height:   float32(cellHeight) * 2,
					velocity: 0,
				}
				continue
			}

			cell.velocity += spreadFactor * laplacian(x, y)
			cell.velocity -= gravity * cell.height
			cell.velocity *= dampening
			cell.height += cell.velocity

			if cell.height < 0 {
				cell.height = 0
				cell.velocity = 0
			}
			maxWater := float32(cellHeight) * 1.5
			if cell.height > maxWater {
				cell.height = maxWater
				cell.velocity *= 0.5
			}

			newWaterGrid[x][y] = cell
		}
	}
	for x := 0; x < gridWidth; x++ {
		for y := 0; y < gridHeight; y++ {
			waterGrid[x][y] = newWaterGrid[x][y]
		}
	}
}

func disturbAt(gridX, gridY int, strength float32) {
	for dx := -3; dx <= 3; dx++ {
		for dy := -3; dy <= 3; dy++ {
			x := gridX + dx
			y := gridY + dy
			if x < 0 || x >= gridWidth || y < 0 || y >= gridHeight {
				continue
			}
			cell := &waterGrid[x][y]
			if isBedrock(cell.height) {
				continue
			}
			dist := float32(math.Sqrt(float64(dx*dx + dy*dy)))
			falloff := 1 - dist/4
			if falloff < 0 {
				continue
			}
			cell.velocity += strength * falloff
			cell.height += strength * 0.5 * falloff
		}
	}
}

func initWater() {
	surfaceY := gridHeight/2 + 1

	for x := 0; x < gridWidth; x++ {
		for y := 0; y < gridHeight; y++ {
			waterGrid[x][y] = WaterCell{}
		}
	}

	for x := 0; x < gridWidth; x++ {
		for y := surfaceY; y < gridHeight-5; y++ {
			waterGrid[x][y].height = float32(cellHeight)
		}
		// Ripple on the surface so waves are visible immediately
		ripple := float32(math.Sin(float64(x)*0.15)) * float32(cellHeight) * 0.4
		waterGrid[x][surfaceY].height += ripple
	}

	for x := 0; x < gridWidth; x++ {
		for y := gridHeight - 5; y < gridHeight; y++ {
			waterGrid[x][y].height = float32(cellHeight) * 2
		}
	}

	// Small random bumps inside the pool
	for i := 0; i < 40; i++ {
		x := rng.Intn(gridWidth)
		y := surfaceY + 1 + rng.Intn(gridHeight-5-surfaceY-1)
		if y < surfaceY+1 {
			continue
		}
		waterGrid[x][y].height += rng.Float32() * float32(cellHeight) * 0.5
	}
}

func drawWater() {
	for x := 0; x < gridWidth; x++ {
		for y := 0; y < gridHeight; y++ {
			cell := waterGrid[x][y]
			drawX := int32(x * cellWidth)
			drawY := int32(y * cellHeight)

			if isBedrock(cell.height) {
				rl.DrawRectangle(drawX, drawY, int32(cellWidth), int32(cellHeight), bedrockColor)
				continue
			}
			if cell.height <= 0 {
				continue
			}

			// Use ceil so small height changes become visible pixels
			level := int32(math.Ceil(float64(cell.height)))
			if level < 1 {
				level = 1
			}
			if level > int32(cellHeight) {
				level = int32(cellHeight)
			}

			alpha := uint8(math.Min(255, 180+float64(cell.height/float32(cellHeight))*75))
			color := rl.NewColor(waterColor.R, waterColor.G, waterColor.B, alpha)

			rl.DrawRectangle(
				drawX,
				drawY+int32(cellHeight)-level,
				int32(cellWidth),
				level,
				color,
			)
		}
	}
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, title)
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	initWater()

	for !rl.WindowShouldClose() {
		updateWater()

		if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			gx := int(rl.GetMouseX()) / cellWidth
			gy := int(rl.GetMouseY()) / cellHeight
			disturbAt(gx, gy, 2.5)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		drawWater()
		rl.DrawText("Clique para criar ondas | ESC para sair", 10, 10, 16, rl.LightGray)
		rl.DrawText(fmt.Sprintf("FPS: %d", rl.GetFPS()), 10, 30, 16, rl.LightGray)
		rl.EndDrawing()
	}
}
