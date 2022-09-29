package util

import (
	"os"
	"strings"
	
	log "github.com/sirupsen/logrus"
)


func InitLogger() *log.Logger {
	Logger := log.New()
	Logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
		ForceColors:  true,
	})

	if os.Getenv("LOG_LEVEL") == "debug" {
		Logger.SetLevel(log.DebugLevel)
	} else {
		Logger.SetLevel(log.InfoLevel)
	}

	if os.Getenv("LOG_FILE") != "" {
		file, err := os.OpenFile(os.Getenv("LOG_FILE"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			Logger.SetOutput(file)
		} else {
			Logger.Info("Failed to log to file, using default stderr")
		}
	}

	return Logger
}

type JError struct {
	Error string `json:"error"`
}

func NewJError(err error) JError {
	jerr := JError{"generic error"}
	if err != nil {
		jerr.Error = err.Error()
	}
	return jerr
}

func NormalizeEmail(email string) string {
	normalizedEmail := strings.ToLower(email)
	normalizedEmail = strings.TrimSpace(normalizedEmail)
	normalizedEmail = strings.Trim(normalizedEmail, "\\")
	return normalizedEmail
}
