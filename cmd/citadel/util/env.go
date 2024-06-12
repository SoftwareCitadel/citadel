package util

import (
	"bufio"
	"os"
	"strings"
)

// RetrieveEnvironmentVariablesFromFile reads a .env file and returns an array of strings in the format "key=value".
func RetrieveEnvironmentVariablesFromFile(path string) ([]string, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var envVariables []string

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || len(strings.TrimSpace(line)) == 0 || !strings.Contains(line, "=") {
			continue
		}

		// Add the line to the array
		envVariables = append(envVariables, line)
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return envVariables, nil
}
