package tools

import (
	"log"

	"go.uber.org/zap"
)

//NewLogger creates new structured zap logger
func NewLogger() *zap.Logger {

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	return logger
}
