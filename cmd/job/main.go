// cmd/api/main.go
package main

import (
	"github.com/caarlos0/env/v11"
	telegram "github.com/go-telegram/bot"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/job"
	"github.com/peatch-io/peatch/internal/notification"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type config struct {
	DatabaseURL   string `env:"DATABASE_URL,required"`
	BotToken      string `env:"BOT_TOKEN,required"`
	ImgServiceURL string `env:"IMG_SERVICE_URL,required"`
	WebAppURL     string `env:"WEB_APP_URL,required"`
	GroupChatID   int64  `env:"GROUP_CHAT_ID,required"`
}

func main() {
	cfg := loadConfig()

	pg, err := db.New(cfg.DatabaseURL)

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	bot, err := telegram.New(cfg.BotToken)

	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	notifier := notification.NewTelegramNotifier(bot)

	notifyJob := job.NewNotifyJob(pg, notifier, cfg.ImgServiceURL, cfg.WebAppURL, cfg.GroupChatID)

	jobs := []*job.Job{
		job.NewJob("NotifyUserReceivedCollaborationRequest", 120*time.Second, notifyJob.NotifyUserReceivedCollaborationRequest),
		job.NewJob("NotifyNewCollaboration", 100*time.Second, notifyJob.NotifyNewCollaboration),
		job.NewJob("NotifyNewUserProfile", 10*time.Second, notifyJob.NotifyNewUserProfile),
		job.NewJob("NotifyCollaborationRequest", 110*time.Second, notifyJob.NotifyCollaborationRequest),
	}

	sc := job.NewScheduler(jobs)
	sc.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
}

func loadConfig() config {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}
	return cfg
}
