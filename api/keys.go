package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// api key and token setup
func ConfigSetup() {

	// api key and token setup
	// Read the configuration from config.json
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return
	}
	defer configFile.Close()

	var config Config
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		fmt.Println("Error decoding config file:", err)
		return
	}

	// Use the paths from the configuration
	clientSecretFile := config.ClientSecretPath

	// Load environment variables from .env file
	if err := loadEnvFromFile(clientSecretFile); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

}

func loadEnvFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			err := os.Setenv(key, value)
			if err != nil {
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
