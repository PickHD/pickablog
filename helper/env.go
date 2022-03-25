package helper

import (
	"os"
	"strconv"
)
// GetEnvString return env string based on key
func GetEnvString(key string) string {
	return os.Getenv(key)
}

// GetEnvString return env integer based on key
func GetEnvInt(key string) int {
	getEnv := os.Getenv(key)
	intEnv,_ := strconv.Atoi(getEnv)

	return intEnv
}