package utils

import (
	"encoding/json"
	"os"
	"regexp"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

var EventNameRegex = regexp.MustCompile(`^[a-zA-Z0-9-_ ]+$`)

func GetEnvOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

var AppendContextError = func(context *gin.Context, err error) {
	parsedErr := context.Error(err)
	log.Logger.Info().Err(parsedErr).Msg("appended context error")
}

var GetJsonStringFromStruct = func(data interface{}) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	jsonString := string(jsonBytes)
	return jsonString, nil
}
