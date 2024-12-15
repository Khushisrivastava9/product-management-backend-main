package services

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger() {
	Logger = logrus.New()
	Logger.Out = os.Stdout
	Logger.SetFormatter(&logrus.JSONFormatter{})
}

func LogRequest(method, path string, statusCode int, duration time.Duration) {
	Logger.WithFields(logrus.Fields{
		"method":      method,
		"path":        path,
		"status_code": statusCode,
		"duration":    duration.Seconds(),
	}).Info("Request processed")
}

func LogAPIError(method, path string, statusCode int, err error) {
	Logger.WithFields(logrus.Fields{
		"method":      method,
		"path":        path,
		"status_code": statusCode,
		"error":       err.Error(),
	}).Error("API error")
}

func LogImageProcessingEvent(event, imageURL string, success bool, err error) {
	fields := logrus.Fields{
		"event":     event,
		"image_url": imageURL,
		"success":   success,
	}
	if err != nil {
		fields["error"] = err.Error()
	}
	Logger.WithFields(fields).Info("Image processing event")
}
