package env

import (
	"os"
	"strconv"
)

func GetEnvString(key, DefaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return DefaultValue
}
func GetEnvInt(key string, DefaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if IntValue, err := strconv.Atoi(value); err == nil {
			return IntValue
		}
	}
	return DefaultValue
}
