package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	// "math/rand"
	"os"
	// "fmt"
	"time"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {

	machine := Machine{}
	// machine.Init("../chip8-test-suite/bin/1-chip8-logo.ch8")
	// machine.Init("../chip8-test-suite/bin/2-ibm-logo.ch8")
	rom_file := os.Args[1]
	machine.Init(rom_file)
	
	defer machine.Deinit()

	DisplayInit()
	defer DisplayDeinit()

	timer_count := time.Now()

	for !rl.WindowShouldClose() {
		err := machine.Clocktick()
		// check(err)
		if err != nil {
			time.Sleep(5 * time.Second)
			panic(err)
		}
		DisplayUpdate(machine.DisplayBuf)
		time.Sleep(2 * time.Millisecond)

		cur_time := time.Now()
		if cur_time.Sub(timer_count).Milliseconds() >= (1000/60){
			machine.TimerUpdate()
			timer_count = time.Now()
		}
	}

	// var instr [2]byte

	// instrloop:
	// for {
	// 	_, err = f.Read(instr[:])
	// 	check(err)
	
	// 	fmt.Printf("instruction: %x : %x\n", instr[0], instr[1])
	
	// 	switch {
	// 	case instr[0] == 0 && instr[1] == 0xE0:
	// 		fmt.Println("Clear screen instruction")
	// 	case instr[0] >> 4 == 6:
	// 		fmt.Println("Load immediate instruction")
	// 	case instr[0] >> 4 == 0xA:
	// 		fmt.Println("Load index instruction")
	// 	case instr[0] >> 4 == 0xD:
	// 		fmt.Println("Draw sprite instruction")
	// 	case instr[0] >> 4 == 0x1:
	// 		fmt.Println("Jump instruction")
	// 	default:
	// 		fmt.Println("Unknown instruction")
	// 		break instrloop
	// 	}
	
	// }
	

	// var display_buffer [32]uint64
	
	// for i := 0; i < 32; i++ {
	// 	display_buffer[i] = rand.Uint64()
	// }	

	// Init()
	// defer Deinit()
	// for !rl.WindowShouldClose() {
	// 	Update(display_buffer)
	// }


}