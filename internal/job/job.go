package job

import (
	"context"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/notification"
	"log"
)

func Run(ctx context.Context, storage *db.Storage, notifier *notification.Notifier) error {
	log.Println("starting job")

	return nil
}
