package main

import (
	"context"
	"corvette/internal/camera"
	"corvette/internal/config"
	"corvette/internal/object_detection"
	"corvette/internal/platform/handler"
	"corvette/internal/streamer"
	"corvette/internal/vendors"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

func pprofStuffs() {
	go func() {
		http.ListenAndServe("localhost:6767", nil)
	}()
}

func main() {
	pprofStuffs()
	coreCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	httpHandler := handler.NewHttpHandler()
	httpHandler.Start(":8080")

	slog.Info("Corvette started.")
	config := config.ReadConfig()

	objectDetectionModel := object_detection.NewObjectDetectionModelInstance(config.OnnxDllPath, config.ObjDetectionModel)

	vendorsFromConfig := vendors.VendorMapper(config.Cameras)
	streamers := streamer.StreamerMapper(vendorsFromConfig)
	cameraRegistry := camera.CreateCameraRegistry(coreCtx, objectDetectionModel)
	cameraRegistry.RegisterArrStreamers(streamers)
	cameraRegistry.StartAllRegisteredCameras()

	<-coreCtx.Done()
	cameraRegistry.WaitToClose()
	slog.Info("Corvette shutting down.")
}
