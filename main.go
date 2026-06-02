package main

import (
	"corvette/internal/configs"
	"corvette/internal/recorder"
)

const (
	width  = 1920
	height = 1080
)

func main() {
	config := configs.ReadConfig()

	recManager := recorder.CreateRecordingManager(config.Cameras)
	recManager.StartAllRecording()

	select {}
}
