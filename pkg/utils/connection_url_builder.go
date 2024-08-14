package utils

import (
	"fmt"
	"os"
)

// ConnectionURLBuilder func for building URL connection.
func ConnectionURLBuilder(n string) (string, error) {
	switch n {
	case "redis":
		return fmt.Sprintf(
			"%s:%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT"),
		), nil
	case "fiber":
		return fmt.Sprintf(
			"%s:%s",
			os.Getenv("APP_HOST"),
			os.Getenv("APP_PORT"),
		), nil
	default:
		return "", fmt.Errorf("connection name '%v' is not supported", n)
	}
}
