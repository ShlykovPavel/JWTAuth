package scheduler

import (
	"context"
	"log/slog"
	"time"
)

type TokenRefresher struct {
	logger     *slog.Logger
	timer      *time.Timer
	cancelFunc context.CancelFunc
}

func NewScheduler(logger *slog.Logger) *TokenRefresher {
	return &TokenRefresher{
		logger: logger,
	}
}

func (t *TokenRefresher) ScheduleRefresh(expiryTime time.Time, refreshFunc func()) {
	const op = "scheduler.scheduleRefresh"
	log := slog.With(
		slog.String("op", op),
		slog.Any("expiryTime", expiryTime))
	// 1. Отменить предыдущий таймер (если был)
	if t.timer != nil {
		log.Info("Timer is already running. Stop actual timer")
		t.timer.Stop()
	}

	// 2. Вычислить время до обновления
	refreshDuration := time.Until(expiryTime) - 1*time.Minute
	log.Info("Calculate time for refresh.", refreshDuration)

	// 3. Запланировать обновление
	t.timer = time.NewTimer(refreshDuration)
	ctx, cancel := context.WithCancel(context.Background())
	t.cancelFunc = cancel
	log.Info("Schedule timer for refresh.", refreshDuration)

	go func() {
		log.Info("Starting timer for refresh.", refreshDuration)
		select {
		case <-t.timer.C:
			t.logger.Info("triggering token refresh")
			refreshFunc()
		case <-ctx.Done():
			t.logger.Debug("refresh canceled")
		}
	}()
}

func (t *TokenRefresher) Stop() {
	if t.cancelFunc != nil {
		t.cancelFunc()
	}
	if t.timer != nil {
		t.timer.Stop()
	}
}
