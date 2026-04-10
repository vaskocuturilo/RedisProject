package utils

import (
	"net"
	"os"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func LoadRedisConfig() *RedisConfig {
	host := getEnv("REDIS_HOST", "localhost")
	port := getEnv("REDIS_PORT", "6379")

	return &RedisConfig{
		Addr:     net.JoinHostPort(host, port),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	}
}
