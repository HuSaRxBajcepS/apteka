package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Server struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"server"`

	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
	} `json:"database"`

	JWTSecret string `json:"jwt_secret"`

	Prescription struct {
		ExpirationDays int `json:"expiration_days"`
	} `json:"prescription"`
}

func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	cfg := &Config{}
	err = json.NewDecoder(
		file,
	).Decode(cfg)
	return cfg, err
}
