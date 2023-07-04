package utils

import (
	"log"

	"go.uber.org/zap"
)

// DefaultLogger is the default logger for this project
var DefaultLogger *zap.Logger

// InitLogger instantiates a logger using Zap
func InitLogger() error {
	var err error
	if DefaultLogger == nil {
		DefaultLogger, err = zap.NewProduction()
		if err != nil {
			log.Fatalf("can't initialize zap logger: %v", err)
		}
	}

	defer DefaultLogger.Sync()
	return nil
}
