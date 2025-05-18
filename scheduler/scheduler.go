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
