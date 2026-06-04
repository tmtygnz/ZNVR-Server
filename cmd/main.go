package main

import (
	"corvette/internal/config"
	"log/slog"
)

func main() {
	slog.Info("Corvette started.")
	_ = config.ReadConfig()
}
