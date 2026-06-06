package camera

import (
	"context"
	"corvette/internal/object_detection"
	"corvette/internal/streamer"
	"log/slog"

	"golang.org/x/sync/errgroup"
)

type OnlineStatus int

const (
	Online = iota
	StartingUp
	Offline
)

type CameraHandler struct {
	streamer    streamer.Streamer
	context     context.Context
	status      OnlineStatus
	objDetModel *object_detection.ObjectDetectionHandler
}

func CreateCameraHandler(streamer streamer.Streamer, ctx context.Context, objDet *object_detection.ObjectDetectionHandler) *CameraHandler {
	return &CameraHandler{
		streamer:    streamer,
		context:     ctx,
		status:      Offline,
		objDetModel: objDet,
	}
}

func (ch *CameraHandler) StartAllFunctions() {
	g, ctx := errgroup.WithContext(ch.context)

	g.Go(func() error {
		return ch.streamer.StartRecording(ctx)
	})
	g.Go(func() error {
		return ch.streamer.StartAIStreaming(ctx)
	})
	g.Go(func() error {
		return ch.StartAIInference(ctx)
	})
	ch.status = Online

	if err := g.Wait(); err != nil {
		slog.Error("Camera function returned", "err", err.Error())
		ch.status = Offline
	}
}

func (ch *CameraHandler) StartAIInference(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			ch.objDetModel.ModelInstance().Destroy()
			slog.Info("Closing object detection inference loop. Other function crashed.")
			return ctx.Err()
		default:
			frame, ok := ch.streamer.GetAIFrame()
			if !ok {
				ch.objDetModel.ModelInstance().Destroy()
				slog.Warn("Get AI Frame is not okay. Closing.")
				return ctx.Err()
			}
			err := ch.objDetModel.DoInference(frame)
			if err != nil {
				return err
			}
		}
	}
}
