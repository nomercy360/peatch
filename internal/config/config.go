package config

import (
	"fmt"
	"github.com/go-playground/validator"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Host            string         `yaml:"host" validate:"required"`
	Port            int            `yaml:"port" validate:"required,gt=0"`
	DBURL           string         `yaml:"mongo_uri" validate:"required"`
	DBName          string         `yaml:"mongo_db" validate:"required"`
	JWTSecret       string         `yaml:"jwt_secret" validate:"required"`
	Telegram        TelegramConfig `yaml:"telegram" validate:"required"`
	LogLevel        string         `yaml:"log_level" validate:"required,oneof=debug info warn error"`
	AWSConfig       AWSConfig      `yaml:"aws" validate:"required"`
	AssetsURL       string         `yaml:"assets_url" validate:"required,url"`
	ImageServiceURL string         `yaml:"image_service_url" validate:"required,url"`
	OpenAIAPIKey    string         `yaml:"openai_api_key" validate:"required"`
	DBPath          string         `yaml:"db_path"`
}

type AWSConfig struct {
	AccessKey string `yaml:"access_key_id" validate:"required"`
	SecretKey string `yaml:"secret_access_key" validate:"required"`
	Bucket    string `yaml:"bucket" validate:"required"`
	Endpoint  string `yaml:"endpoint" validate:"required"`
}

type TelegramConfig struct {
	BotToken         string `yaml:"bot_token" validate:"required"`
	AdminBotToken    string `yaml:"admin_bot_token" validate:"required"`
	AdminChatID      int64  `yaml:"admin_chat_id" validate:"required"`
	CommunityChatID  int64  `yaml:"community_chat_id" validate:"required"`
	WebAppURL        string `yaml:"webapp_url" validate:"required,url"`
	BotWebApp        string `yaml:"bot_webapp" validate:"required"`
	WebhookURL       string `yaml:"webhook_url" validate:"required,url"`
	TestNotification bool   `yaml:"test_notification"`
}

func LoadConfig() (*Config, error) {
	configFilePath := "config.yml"
	if envPath := os.Getenv("CONFIG_FILE_PATH"); envPath != "" {
		configFilePath = envPath
	}

	file, err := os.Open(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("configuration validation error: %w", err)
	}

	return &cfg, nil
}
