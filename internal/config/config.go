package config

import "os"

type Config struct {
	Host        string
	Port        int
	DatabaseURL string
	RedisURL    string
	HostKeyPath string
}

func Load() *Config {
	return &Config{
		Host:        envOr("SSH_HOST", "0.0.0.0"),
		Port:        2222,
		DatabaseURL: envOr("DATABASE_URL", "postgres://packets:packets@localhost:5432/packets?sslmode=disable"),
		RedisURL:    envOr("REDIS_URL", "redis://localhost:6379/0"),
		HostKeyPath: envOr("HOST_KEY_PATH", ".ssh/host_key"),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
