package config

import "os"

type Config struct {
	Port     string
	DataPath string
}

func Load() *Config {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080" // Default port
	}
	dataPath := os.Getenv("DATA_PATH")
	if dataPath == "" {
		dataPath = "data" // Default data path
	}

	return &Config{
		Port:     port,
		DataPath: dataPath,
	}
}
