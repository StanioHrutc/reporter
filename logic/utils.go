package logic

import (
	"encoding/json"
	"fmt"
	"os"
)

const ConfigPath = "etc/config.json"

type Config struct {
	Queue       QueueConfig        `json:"queue"`
	BoardConfig BoardServiceConfig `json:"status_board"`
}

type BoardServiceConfig struct {
	StatusBoardHost string `json:"status_board_service_host"`
	StatusBoardPort int    `json:"status_board_service_port"`
}

type QueueConfig struct {
	Url         string `json:"url"`
	QueueName   string `json:"name"`
	RetryAmount int    `json:"retry_amount"`
}

func GetConfig() *Config {
	configFile, _ := os.Open(ConfigPath)
	defer configFile.Close() // defer the closing for safety purposes

	decoder := json.NewDecoder(configFile)
	configuration := &Config{}

	err := decoder.Decode(configuration)
	if err != nil {
		fmt.Printf("Got error while reading configuration file: %v. \nExiting", err)
		os.Exit(1)
	}

	return configuration
}
