package object_detection

import (
	"log/slog"
	"strings"
)

type ObjectDetectionHandler struct {
	modelInstance *ObjectDetectionModelInstance
	lastFrame     []float32
}

func CreateObjectDetectionHandler(modelInstance *ObjectDetectionModelInstance) *ObjectDetectionHandler {
	return &ObjectDetectionHandler{
		modelInstance: modelInstance,
	}
}

// Process frames and returns detection that is within the set threshold.
// The `frame` parameter must already be in a flat non-interleaved [R,R,R,...,G,G,G,...,B,B,B] Format
func (odh *ObjectDetectionHandler) DoInference(frame []float32) error {
	copy(odh.modelInstance.InputTensor.GetData(), frame)
	err := odh.modelInstance.Session.Run()
	if err != nil {
		return err
	}

	data := odh.modelInstance.OutputTensor.GetData()
	for i := range 300 {
		offset := i * 6
		score := data[offset+4]

		if score > 0.4 {
			classId := int(data[offset+5])
			classStr := strings.TrimSpace(odh.modelInstance.Categories[classId])

			// Extracting coordinates for the structured log
			bbox := []float32{
				data[offset+0], // x1
				data[offset+1], // y1
				data[offset+2], // x2
				data[offset+3], // y2
			}

			slog.Info("Detection found",
				"class", classStr,
				"score", score,
				"position", bbox,
			)
		}
	}
	odh.lastFrame = frame

	return nil
}

func (odh *ObjectDetectionHandler) ModelInstance() *ObjectDetectionModelInstance {
	return odh.modelInstance
}
