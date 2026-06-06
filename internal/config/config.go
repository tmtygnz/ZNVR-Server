package config

type Config struct {
	AiScalingSize     int          `toml:"ai_scaling_size"`
	Cameras           []CameraInfo `toml:"cameras"`
	OnnxDllPath       string       `toml:"onnx_dll_path"`
	ObjDetectionModel string       `toml:"obj_detection_model"`
}

type CameraInfo struct {
	URL     string `toml:"url"`
	SURL    string `toml:"sub_url"`
	Type    string `toml:"type"`
	CamName string `toml:"cam_name"`
}
