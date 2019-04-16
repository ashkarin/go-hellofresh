package config

import (
	"encoding/json"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// DBConfig database config
type DBConfig struct {
	Server   string `json:"server"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

// Config is Server and DB configuration
type Config struct {
	DB      DBConfig `json:"db"`
	Address string   `json:"address"`
	Port    string   `json:"port"`
	Timeout int      `json:"timeout"`
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

// GetConfigFromEnv get config from environment variables
func GetConfigFromEnv() (*Config, error) {
	timeout, err := strconv.ParseInt(getenv("SRV_TIMEOUT", "10"), 10, 32)
	if err != nil {
		log.Error("Wrong value on SRV_TIMEOUT. It will be set to 10.")
		timeout = 10
	}
	cfg := &Config{
		DB: DBConfig{
			Server:   getenv("DB_HOST", "mongodb"),
			Port:     getenv("DB_PORT", "27017"),
			Username: getenv("DB_USER", ""),
			Password: getenv("DB_PASS", ""),
			DBName:   getenv("DB_NAME", "hellofresh"),
		},
		Address: getenv("SRV_HOST", ""),
		Port:    getenv("SRV_PORT", "8080"),
		Timeout: int(timeout),
	}
	return cfg, nil
}

// GetConfig get config from the JSON file
func GetConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	config := &Config{}
	dec := json.NewDecoder(file)
	if err := dec.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
