package env

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func MustString(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("required ENV %q is not set", key)
	}
	if value == "" {
		log.Fatalf("required ENV %q is empty", key)
	}
	return value
}

func MustInt(key string) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("required ENV %q is not set", key)
	}
	if value == "" {
		log.Fatalf("required ENV %q is empty", key)
	}
	res, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		log.Fatalf("required ENV %q must be a number but it's %q", key, value)
	}
	return int(res)
}

func MustBool(key string) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("required ENV %q is not set", key)
	}
	if value == "true" || value == "1" {
		return true
	}
	return false
}

func LoadEnvFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		os.Setenv(key, value)
	}

	return scanner.Err()
}
