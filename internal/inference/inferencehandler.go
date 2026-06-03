package inference

import (
	"bytes"
	"corvette/internal/capture"
	"log"
	"strings"
)

type InferenceHandler struct {
	capturer      capture.Capturer
	modelInstance *ModelInstance
	lastFrame     []byte
}

func CreateNewInferenceHanlder(capturer capture.Capturer, modelInstance *ModelInstance) *InferenceHandler {
	return &InferenceHandler{
		capturer,
		modelInstance,
		nil,
	}
}

func (ih *InferenceHandler) StartHandler() error {
	err := ih.capturer.StartAIStreamer()
	if err != nil {
		return err
	}

	go ih.inferenceLoop()

	return nil
}

func (ih *InferenceHandler) inferenceLoop() {
	for {
		frame, ok := ih.capturer.GetCurrentAIFrame()
		if !ok {
			log.Printf("AI Stream stopped, stopping inference loop")
			break
		}
		if bytes.Equal(frame, ih.lastFrame) {
			continue
		}

		deltaFrame := ih.diffFrame(frame)
		log.Printf("%f", deltaFrame)

		if deltaFrame < 0.07 {
			ih.lastFrame = frame
			log.Println("Skipping yolo, not that much changed.")
			continue
		}

		dataTensor := ih.prepareInput(frame)

		copy(ih.modelInstance.InputTensor.GetData(), dataTensor)

		ih.modelInstance.Session.Run()

		data := ih.modelInstance.OutputTensor.GetData()
		for i := range 300 {
			offset := i * 6
			score := data[offset+4]
			if score > 0.6 {
				classId := int(data[offset+5])
				classStr := strings.TrimSpace(ih.modelInstance.Categories[classId])
				log.Printf("Found class %s with score %f", classStr, score)
			}
		}
		ih.lastFrame = frame
	}
}

func (ih *InferenceHandler) prepareInput(frame []byte) []float32 {
	const (
		width  = 640
		height = 640
		size   = width * height
	)

	input := make([]float32, size*3)

	for i := range size {
		// RGB Interleaved (byte) -> RGB Planar (float32)
		input[i] = float32(frame[i*3]) / 255.0          // R
		input[i+size] = float32(frame[i*3+1]) / 255.0   // G
		input[i+size*2] = float32(frame[i*3+2]) / 255.0 // B
	}

	return input
}

func (ih *InferenceHandler) diffFrame(newf []byte) float64 {
	const totalPixels = 640 * 640
	const noiseThreshold = 8
	changedPixels := 0
	base := ih.lastFrame

	for i := 0; i < len(base); i += 3 {
		rDif := int(base[i]) - int(newf[i])
		if rDif < 0 {
			rDif = -rDif
		}

		gDif := int(base[i+1]) - int(newf[i+1])
		if gDif < 0 {
			gDif = -gDif
		}

		bDif := int(base[i+2]) - int(newf[i+2])
		if bDif < 0 {
			bDif = -bDif
		}

		if rDif > noiseThreshold || gDif > noiseThreshold || bDif > noiseThreshold {
			changedPixels++
		}
	}

	return float64(changedPixels) / float64(totalPixels)
}
