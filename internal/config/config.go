package config

import "os"

const (
	DEFAULT_CONFIG_PATH = "/data/config.json"
	TEMP_PATH           = "/temp"
)

func GetConfigPath() string {
	e := os.Getenv("CONFIG_PATH")
	if e == "" {
		return DEFAULT_CONFIG_PATH
	}
	return e
}

func GetTempPath() string {
	e := os.Getenv("TEMP_PATH")
	if e == "" {
		return TEMP_PATH
	}
	return e
}
