package config

import "os"

type Config struct {
	Address string
	DB      string
	DBPort  string
}

func GetConfig() Config {
	return Config{
		Address: os.Getenv("RATE_LIMITER_SERVICE_ADDRESS"),
		DB:      os.Getenv("DB"),
		DBPort:  os.Getenv("DBPORT"),
	}
}
