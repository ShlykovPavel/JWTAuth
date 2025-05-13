package scheduler

import (
	"context"
	"log/slog"
	"time"
)

type Scheduler struct {
	timer      *time.Timer
	cancelFunc context.CancelFunc
	onRefresh  func()
	logger     *slog.Logger
}

func NewScheduler(onRefresh func(), logger *slog.Logger) *Scheduler {
	return &Scheduler{
		onRefresh: onRefresh,
		logger:    logger,
	}
}

//func (s *Scheduler) ScheduleRefresh(expiryTime time.Time, refreshFunc func()) {
//	const op = "scheduler.scheduleRefresh"
//	log := s.logger.With(
//		slog.String("op", op),
//		slog.Any("expiryTime", expiryTime))
//	// 1. Отменить предыдущий таймер (если был)
//	if s.timer != nil {
//		log.Info("Timer is already running. Stop actual timer")
//		s.timer.Stop()
//	}
//
//	// 2. Вычислить время до обновления
//	refreshDuration := time.Until(expiryTime) - 1*time.Minute
//	log.Info("Calculate time for refresh.", refreshDuration)
//
//	// 3. Запланировать обновление
//	s.timer = time.NewTimer(refreshDuration)
//	ctx, cancel := context.WithCancel(context.Background())
//	s.cancelFunc = cancel
//	log.Info("scheduled token refresh",
//		"refresh_in", refreshDuration.String(),
//		"refresh_at", time.Now().Add(refreshDuration).Format(time.RFC3339))
//
//	go func() {
//		defer log.Debug("refresh goroutine stopped")
//		log.Info("Starting timer for refresh.", refreshDuration)
//		select {
//		case <-s.timer.C:
//			s.logger.Info("triggering token refresh")
//			refreshFunc()
//		case <-ctx.Done():
//			s.logger.Debug("refresh canceled")
//		}
//	}()
//}

func (s *Scheduler) Stop() {
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
	if s.timer != nil {
		s.timer.Stop()
	}
}

func (s *Scheduler) ScheduleRefresh(expiry time.Time) {
	const op = "scheduler.scheduleRefresh"
	log := s.logger.With(
		slog.String("op", op))
	s.Stop()

	refreshIn := time.Until(expiry) - 1*time.Minute
	log.Debug("Calculating time for init refresh: ", refreshIn)
	if refreshIn < 10*time.Second {
		refreshIn = 10 * time.Second
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel

	s.timer = time.NewTimer(refreshIn)
	go func() {
		select {
		case <-s.timer.C:
			s.onRefresh()
		case <-ctx.Done():
			s.logger.Debug("refresh canceled")
		}
	}()
}
