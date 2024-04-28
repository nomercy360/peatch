package main

import (
	"github.com/caarlos0/env/v11"
	_ "github.com/peatch-io/peatch/docs"
	"github.com/peatch-io/peatch/internal/bot"
	"github.com/peatch-io/peatch/internal/db"
	storage "github.com/peatch-io/peatch/internal/s3"
	"log"
	"net/http"
)

type config struct {
	DatabaseURL string `env:"DATABASE_URL,required"`
	Server      ServerConfig
	BotToken    string `env:"BOT_TOKEN,required"`
	AWS         AWSConfig
	ExternalURL string `env:"EXTERNAL_URL,required"`
	WebAppURL   string `env:"WEB_APP_URL,required"`
}

type ServerConfig struct {
	Port string `env:"SERVER_PORT" envDefault:"8080"`
	Host string `env:"SERVER_HOST" envDefault:"localhost"`
}

type AWSConfig struct {
	AccessKey string `env:"AWS_ACCESS_KEY_ID,required"`
	SecretKey string `env:"AWS_SECRET_ACCESS_KEY,required"`
	Bucket    string `env:"AWS_BUCKET,required"`
	Endpoint  string `env:"AWS_ENDPOINT,required"`
}

func main() {
	cfg := loadConfig()
	dbClient, err := db.New(cfg.DatabaseURL)

	if err != nil {
		log.Fatalf("Failed to initialize database: %v\n", err)
	}

	defer dbClient.Close()

	s3Client, err := storage.NewS3Client(cfg.AWS.AccessKey, cfg.AWS.SecretKey, cfg.AWS.Endpoint, cfg.AWS.Bucket)

	if err != nil {
		log.Fatalf("Failed to initialize S3 client: %v", err)
	}

	botCfg := bot.Config{
		BotToken:    cfg.BotToken,
		WebAppURL:   cfg.WebAppURL,
		ExternalURL: cfg.ExternalURL,
	}

	tgBot := bot.New(dbClient, s3Client, botCfg)

	if err := tgBot.SetWebhook(); err != nil {
		log.Fatalf("Failed to set webhook: %v", err)
	}

	log.Printf("Webhook set to %s", cfg.ExternalURL+"/webhook")

	http.HandleFunc("/webhook", tgBot.HandleWebhook)

	log.Printf("Starting server at %s:%s...", cfg.Server.Host, cfg.Server.Port)

	if err := http.ListenAndServe(cfg.Server.Host+":"+cfg.Server.Port, nil); err != nil {
		log.Fatalf("Failed to start webhook server: %v", err)
	}
}

func loadConfig() config {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}
	return cfg
}
