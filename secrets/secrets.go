package secrets

import (
	"os"
	"strings"
)

// ReadSecret reads a file from the given path and trims whitespace
func ReadSecret(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}
