package util

import "os"

func GetEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return value
}

func DefaultStr(value, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}
