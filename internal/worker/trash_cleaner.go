package worker

import (
	"context"
	"time"

	"go.uber.org/zap"

	"CloudBox/internal/usecase/trash"
)

const cleanerBatch = 100

func RunTrashCleaner(ctx context.Context, uc *trash.Usecase, l *zap.Logger, every time.Duration) {
	run := func() {
		files, folders, err := uc.PurgeExpired(ctx, cleanerBatch)
		if err != nil {
			l.Error("trash cleaner", zap.Error(err))
			return
		}
		if files > 0 || folders > 0 {
			l.Info("trash cleaner purged", zap.Int("files", files), zap.Int("folders", folders))
		}
	}

	run()

	ticker := time.NewTicker(every)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			run()
		}
	}
}
