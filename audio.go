package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"		
)

var (
	music rl.Music
)

func AudioInit(){
	rl.InitAudioDevice()

	music = rl.LoadMusicStream("music.mp3")
	
}

func AudioPlay(){
	rl.PlayMusicStream(music)
	rl.UpdateMusicStream(music)
}

func AudioStop(){
	rl.StopMusicStream(music)
}

func AudioUpdate(){
	rl.UpdateMusicStream(music)
}

func AudioDeinit(){
	rl.UnloadMusicStream(music)	
	rl.CloseAudioDevice() // Close audio device (music streaming is automatically stopped)
}
