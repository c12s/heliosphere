package config

import "os"

type Config struct {
	Address                   string
	RateLimiterServiceAddress string
}

func GetConfig() Config {
	return Config{
		RateLimiterServiceAddress: os.Getenv("RATE_LIMITER_SERVICE_ADDRESS"),
		Address:                   os.Getenv("GATEWAY_ADDRESS"),
	}
}
