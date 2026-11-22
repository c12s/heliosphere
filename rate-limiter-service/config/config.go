package config

import (
	"os"
	"strconv"
)

type Config struct {
	DB     string
	DBPort string
	// limit of the rate-limiter auto-created for a method's first call
	DefaultReqPerSec int64
	DefaultBurst     int64
}

func GetConfig() Config {
	return Config{
		DB:               os.Getenv("DB_CONSUL"),
		DBPort:           os.Getenv("DBPORT_CONSUL"),
		DefaultReqPerSec: envInt("RL_DEFAULT_REQ_PER_SEC", 1),
		DefaultBurst:     envInt("RL_DEFAULT_BURST", 1),
	}
}

func envInt(key string, def int64) int64 {
	v, err := strconv.ParseInt(os.Getenv(key), 10, 64)
	if err != nil || v < 1 {
		return def
	}
	return v
}
