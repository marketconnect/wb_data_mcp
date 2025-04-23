package config

import (
	"github.com/caarlos0/env/v6"
)

// Config holds all the application configuration
type Config struct {
	Server struct {
		IP   string `env:"WB_DATA_MCP_IP" envDefault:"0.0.0.0"`
		Port string `env:"WB_DATA_MCP_PORT" envDefault:"8082"` // Changed default port
	}

	Clickhouse struct {
		Username string `env:"CH_USERNAME" envDefault:"default"`
		Password string `env:"CLICKHOUSE_PASSWORD" envRequired:"true"`
		Database string `env:"CLICKHOUSE_DATABASE" envRequired:"true"`
		Host     string `env:"CH_HOST" envDefault:"localhost"`
		Port     string `env:"CH_PORT" envDefault:"9000"`
	}

	Postgres struct {
		Username string `env:"POSTGRES_USERNAME" envDefault:"postgres"`
		Password string `env:"CLICKHOUSE_PASSWORD" envRequired:"true"` // Reuse same password
		Host     string `env:"PG_HOST" envDefault:"localhost"`
		Port     string `env:"PG_PORT" envDefault:"5432"`
		Database string `env:"CLICKHOUSE_DATABASE" envRequired:"true"` // Reuse same database
	}

	Redis struct {
		Host     string `env:"REDIS_HOST" envDefault:"localhost"`
		Port     string `env:"REDIS_PORT" envDefault:"6379"`
		Password string `env:"REDIS_PASSWORD" envDefault:""`
		DB       int    `env:"REDIS_DB" envDefault:"0"`
	}

	Telegram struct {
		ChatID   string `env:"TELEGRAM_CHAT_ID" envRequired:"true"`
		BotToken string `env:"TELEGRAM_BOT_TOKEN" envRequired:"true"`
	}
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
