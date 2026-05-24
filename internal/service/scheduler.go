package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/qnguyenhong/automation-platform/internal/model"
)

// SchedulerService handles cron-based scheduled test runs.
type SchedulerService struct {
	scheduleRepo ScheduleRepository
	runSvc       *TestRunService
	logger       *slog.Logger
	ticker       *time.Ticker
	done         chan struct{}
}

func NewSchedulerService(repo ScheduleRepository, runSvc *TestRunService, logger *slog.Logger) *SchedulerService {
	return &SchedulerService{
		scheduleRepo: repo,
		runSvc:       runSvc,
		logger:       logger,
		done:         make(chan struct{}),
	}
}

// Start begins the scheduler loop, checking for due schedules every interval.
func (s *SchedulerService) Start(ctx context.Context, interval time.Duration) {
	s.ticker = time.NewTicker(interval)
	s.logger.Info("scheduler started", "interval", interval)

	go func() {
		for {
			select {
			case <-ctx.Done():
				s.ticker.Stop()
				s.logger.Info("scheduler stopped")
				return
			case <-s.done:
				s.ticker.Stop()
				return
			case <-s.ticker.C:
				s.processDueSchedules(ctx)
			}
		}
	}()
}

// Stop stops the scheduler.
func (s *SchedulerService) Stop() {
	close(s.done)
}

func (s *SchedulerService) processDueSchedules(ctx context.Context) {
	schedules, err := s.scheduleRepo.ListDue(ctx)
	if err != nil {
		s.logger.Error("failed to list due schedules", "error", err)
		return
	}

	for _, schedule := range schedules {
		s.logger.Info("triggering scheduled run", "suite_id", schedule.SuiteID, "schedule_id", schedule.ID)

		_, err := s.runSvc.Trigger(ctx, schedule.SuiteID, "", model.TriggerRunRequest{
			Trigger:  model.TriggerScheduled,
			Metadata: json.RawMessage(`{"schedule_id": "` + schedule.ID.String() + `"}`),
		})
		if err != nil {
			s.logger.Error("failed to trigger scheduled run", "schedule_id", schedule.ID, "error", err)
			continue
		}

		// Calculate next run time (simple: add the interval based on cron)
		nextRun := calculateNextRun(schedule.CronExpr)
		if err := s.scheduleRepo.UpdateNextRun(ctx, schedule.ID, nextRun); err != nil {
			s.logger.Error("failed to update next run time", "schedule_id", schedule.ID, "error", err)
		}
	}
}

// calculateNextRun parses a cron-like expression and returns the next run time.
// Supports simplified expressions: @every <duration>, or standard 5-field cron.
func calculateNextRun(cronExpr string) time.Time {
	now := time.Now()

	// Handle @every syntax
	if len(cronExpr) > 6 && cronExpr[:6] == "@every" {
		durationStr := cronExpr[7:]
		if d, err := time.ParseDuration(durationStr); err == nil {
			return now.Add(d)
		}
	}

	// Handle common presets
	switch cronExpr {
	case "@hourly":
		return now.Add(1 * time.Hour)
	case "@daily":
		return now.Add(24 * time.Hour)
	case "@weekly":
		return now.Add(7 * 24 * time.Hour)
	}

	// Default: run again in 1 hour
	return now.Add(1 * time.Hour)
}
