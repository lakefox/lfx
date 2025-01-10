package utils

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

type Post struct {
	ID        int
	Title     string
	URL       string
	Text      string
	Score     int
	Timestamp time.Time
	Username  string
}

type Comment struct {
	PostID    int
	Username  string
	Content   string
	Timestamp time.Time
}

// loadEnvFile manually loads a .env file into the environment variables
func ENV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignore empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("Invalid line in .env file: %s", line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove surrounding quotes, if any
		value = strings.Trim(value, `"'`)

		// Set the environment variable
		os.Setenv(key, value)
	}

	return scanner.Err()
}
