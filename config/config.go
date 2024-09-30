package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type HTTPServerConfig struct {
	Port int `env:"APP_PORT" envDefault:"80"`
}

type LoggerConfig struct {
	Level    string `env:"LOG_LEVEL" envDefault:"debug"`
	Mode     string `env:"LOG_MODE" envDefault:"development"`
	Encoding string `env:"LOG_ENCODING" envDefault:"console"`
}

type MongoConfig struct {
	Database string `env:"MONGODB_DATABASE"`
	URI      string `env:"MONGODB_URI"`
}

type HcmusConfig struct {
	Username string `env:"HCMUS_USERNAME"`
	Password string `env:"HCMUS_PASSWORD"`
	URL      string `env:"HCMUS_URL"`
	SessKey  string `env:"HCMUS_SESSKEY"`
	Cookies  string `env:"HCMUS_COOKIES"`
}

type TelegramConfig struct {
	BotToken string `env:"TELEGRAM_BOT_TOKEN"`
	ChatID   int64  `env:"TELEGRAM_CHAT_ID"`
}

type Config struct {
	HTTPServer     HTTPServerConfig
	Logger         LoggerConfig
	Mongo          MongoConfig
	HcmusConfig    HcmusConfig
	TelegramConfig TelegramConfig
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	cfg := &Config{}
	err = env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
