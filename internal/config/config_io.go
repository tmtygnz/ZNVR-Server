package config

import (
	"log/slog"
	"os"

	"github.com/BurntSushi/toml"
)

func ReadConfig() *Config {
	filePath := "./conf.toml"
	content, err := os.ReadFile(filePath)
	if err != nil {
		slog.Error("Failed to open conf.toml file.", "err", err.Error())
		panic(err)
	}

	var configuration Config
	_, err = toml.Decode(string(content), &configuration)
	if err != nil {
		slog.Error("Failed to decode conf.toml file", "err", err.Error())
		panic(err)
	}

	if err := validate(&configuration); err != nil {
		slog.Error("Config failed validation checks", "cause", err.Error())
		panic(err)
	}

	return &configuration
}

func validate(config *Config) error {
	// Check ai scaling
	if config.AiScalingSize == 0 {
		return MissingFieldError("Check AIScalingSize must exist & not be 0.")
	}

	if config.ObjDetectionModel == "" {
		slog.Warn("`obj_detection_model` is not set. It is okay if you will not be using any AI features.")
	}

	if config.OnnxDllPath == "" {
		slog.Warn("`onnx_dll_path` is not set. It is okay if you will not be using any AI features.")
	}
	return nil
}
