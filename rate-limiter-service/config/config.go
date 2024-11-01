package config

import "os"

type Config struct {
	DB     string
	DBPort string
}

func GetConfig() Config {
	return Config{
		DB:     os.Getenv("DB_CONSUL"),
		DBPort: os.Getenv("DBPORT_CONSUL"),
	}
}
