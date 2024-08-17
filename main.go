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

func machineThread(machine *Machine){
	timer_count := time.Now()
	for {
		err := machine.Clocktick()
		// check(err)
		if err != nil {
			time.Sleep(5 * time.Second)
			panic(err)
		}
		
		time.Sleep(2 * time.Millisecond)
		
		cur_time := time.Now()
		if cur_time.Sub(timer_count).Milliseconds() >= (1000/120){
			machine.TimerUpdate()
			timer_count = time.Now()
		}		
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

	go machineThread(&machine)

	for !rl.WindowShouldClose() {
		DisplayUpdate(machine.DisplayBuf)
	}
}