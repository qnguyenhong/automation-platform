package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Auth     AuthConfig     `yaml:"auth"`
	Worker   WorkerConfig   `yaml:"worker"`
	Log      LogConfig      `yaml:"log"`
}

type ServerConfig struct {
	Host         string        `yaml:"host" env:"SERVER_HOST" default:"0.0.0.0"`
	Port         int           `yaml:"port" env:"SERVER_PORT" default:"8080"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT" default:"30s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT" default:"30s"`
}

type DatabaseConfig struct {
	URL string `yaml:"url" env:"DATABASE_URL" default:"postgres://postgres:postgres@localhost:5432/automation_platform?sslmode=disable"`
}

type AuthConfig struct {
	JWTSecret string        `yaml:"jwt_secret" env:"JWT_SECRET" default:"change-me"`
	JWTExpiry time.Duration `yaml:"jwt_expiry" env:"JWT_EXPIRY" default:"24h"`
}

type WorkerConfig struct {
	ServerURL         string        `yaml:"server_url" env:"WORKER_SERVER_URL" default:"http://localhost:8080"`
	Name              string        `yaml:"name" env:"WORKER_NAME" default:"worker-1"`
	ExecutorTypes     []string      `yaml:"executor_types" env:"WORKER_EXECUTOR_TYPES" default:"api"`
	MaxConcurrent     int           `yaml:"max_concurrent" env:"WORKER_MAX_CONCURRENT" default:"4"`
	PollInterval      time.Duration `yaml:"poll_interval" env:"WORKER_POLL_INTERVAL" default:"2s"`
	HeartbeatInterval time.Duration `yaml:"heartbeat_interval" env:"WORKER_HEARTBEAT_INTERVAL" default:"10s"`
}

type LogConfig struct {
	Level  string `yaml:"level" env:"LOG_LEVEL" default:"info"`
	Format string `yaml:"format" env:"LOG_FORMAT" default:"json"`
}

// Load reads configuration from a YAML file and overrides with environment variables.
func Load(path string) (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host:         "0.0.0.0",
			Port:         8080,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		Database: DatabaseConfig{
			URL: "postgres://postgres:postgres@localhost:5432/automation_platform?sslmode=disable",
		},
		Auth: AuthConfig{
			JWTSecret: "change-me",
			JWTExpiry: 24 * time.Hour,
		},
		Worker: WorkerConfig{
			ServerURL:         "http://localhost:8080",
			Name:              "worker-1",
			ExecutorTypes:     []string{"api"},
			MaxConcurrent:     4,
			PollInterval:      2 * time.Second,
			HeartbeatInterval: 10 * time.Second,
		},
		Log: LogConfig{
			Level:  "info",
			Format: "json",
		},
	}

	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read config file: %w", err)
		}
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("parse config file: %w", err)
		}
	}

	// Override with environment variables
	overrideFromEnv(cfg)

	return cfg, nil
}

func overrideFromEnv(cfg *Config) {
	if v := os.Getenv("SERVER_HOST"); v != "" {
		cfg.Server.Host = v
	}
	if v := os.Getenv("SERVER_PORT"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.Server.Port)
	}
	if v := os.Getenv("DATABASE_URL"); v != "" {
		cfg.Database.URL = v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.Auth.JWTSecret = v
	}
	if v := os.Getenv("JWT_EXPIRY"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Auth.JWTExpiry = d
		}
	}
	if v := os.Getenv("WORKER_SERVER_URL"); v != "" {
		cfg.Worker.ServerURL = v
	}
	if v := os.Getenv("WORKER_NAME"); v != "" {
		cfg.Worker.Name = v
	}
	if v := os.Getenv("WORKER_EXECUTOR_TYPES"); v != "" {
		cfg.Worker.ExecutorTypes = splitAndTrim(v)
	}
	if v := os.Getenv("WORKER_MAX_CONCURRENT"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.Worker.MaxConcurrent)
	}
	if v := os.Getenv("WORKER_POLL_INTERVAL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Worker.PollInterval = d
		}
	}
	if v := os.Getenv("WORKER_HEARTBEAT_INTERVAL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Worker.HeartbeatInterval = d
		}
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.Log.Level = v
	}
	if v := os.Getenv("LOG_FORMAT"); v != "" {
		cfg.Log.Format = v
	}
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// Addr returns the server address in host:port format.
func (c *ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
