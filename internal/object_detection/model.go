package object_detection

import (
	"log"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/yalue/onnxruntime_go"
)

type ObjectDetectionModelInstance struct {
	Session      *onnxruntime_go.AdvancedSession
	InputTensor  *onnxruntime_go.Tensor[float32]
	OutputTensor *onnxruntime_go.Tensor[float32]
	Categories   []string
	mu           sync.Mutex
}

func NewObjectDetectionModelInstance(dllPath string, onnxModelPath string) *ObjectDetectionModelInstance {
	slog.Info("Using onnxlib", "path", dllPath)

	onnxruntime_go.SetSharedLibraryPath(dllPath)

	err := onnxruntime_go.InitializeEnvironment()
	if err != nil {
		slog.Error("Failed to initialize onnx env.", "err", err)
		panic(err)
	}

	inputShape := onnxruntime_go.NewShape(1, 3, 640, 640)
	inputTensor, _ := onnxruntime_go.NewEmptyTensor[float32](inputShape)

	outputShape := onnxruntime_go.NewShape(1, 300, 6)
	outputTensor, _ := onnxruntime_go.NewEmptyTensor[float32](outputShape)

	log.Printf("Using %s for inference.", onnxModelPath)

	session, err := onnxruntime_go.NewAdvancedSession(
		onnxModelPath,
		[]string{"images"},  // Input layer name
		[]string{"output0"}, // Output layer name
		[]onnxruntime_go.Value{inputTensor},
		[]onnxruntime_go.Value{outputTensor},
		nil,
	)
	if err != nil {
		slog.Error("Failed to load onnx model", "err", err)
		panic(err)
	}

	categories := ReadCategories("./models/yolo26n.txt")

	return &ObjectDetectionModelInstance{
		Session:      session,
		InputTensor:  inputTensor,
		OutputTensor: outputTensor,
		Categories:   categories,
	}
}

func (odmi *ObjectDetectionModelInstance) Destroy() error {
	odmi.mu.Lock()
	defer odmi.mu.Unlock()
	if odmi.Session == nil {
		slog.Info("AI model instance already destroyed.")
		return nil
	}
	err := odmi.Session.Destroy()
	odmi.Session = nil
	return err
}

func ReadCategories(filePath string) []string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	dataStr := string(data)
	arr := strings.Split(dataStr, "\n")
	return arr
}
