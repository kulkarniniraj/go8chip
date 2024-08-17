package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"fmt"
	"time"
)

const GridSize int32 = 20

func DisplayInit(){
	rl.InitWindow(64 * GridSize, 32 * GridSize, "raylib [core] example - basic window")
	rl.SetTargetFPS(60)
}

// display_buffer [][]u8
func DisplayUpdate(buffer [32]uint64) {

		var row int32
		var col int32
		start := time.Now()
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		buf2 := buffer

		for row = 0; row < 32; row++ {
			for col = 0; col < 64; col++ {
				if buf2[row] & (1 << 63 >> col) != 0 {
					rl.DrawRectangle(col * GridSize, row * GridSize, GridSize - 2, 
						GridSize - 2, rl.Red)
				}				
			}
		}

		// rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.Green)

		rl.EndDrawing()
		
		// time.Sleep(time.Millisecond * 60)
		end := time.Now()
		fmt.Println("Display time: ", end.Sub(start).Milliseconds())

	
}

func DisplayDeinit(){
	rl.CloseWindow()
}