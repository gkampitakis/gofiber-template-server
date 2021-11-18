package config

import "os"

func GetEnv(value string, defaultValue ...string) string {
	envValue := os.Getenv(value)
	if envValue == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return envValue
}
