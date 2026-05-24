package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/qnguyenhong/automation-platform/internal/model"
)

type DashboardRepository interface {
	GetSummary(ctx context.Context, projectID *uuid.UUID, since time.Time) (*DashboardSummary, error)
	GetRecentRuns(ctx context.Context, projectID *uuid.UUID, limit int) ([]model.TestRun, error)
	GetFlakyTests(ctx context.Context, projectID *uuid.UUID, since time.Time, limit int) ([]FlakyTest, error)
	GetTrends(ctx context.Context, projectID *uuid.UUID, from, to time.Time) ([]TrendPoint, error)
}

type DashboardSummary struct {
	TotalRuns     int64 `json:"total_runs"`
	TotalPassed   int64 `json:"total_passed"`
	TotalFailed   int64 `json:"total_failed"`
	ActiveWorkers int64 `json:"active_workers"`
	PassRate      float64 `json:"pass_rate"`
}

type FlakyTest struct {
	CaseID   uuid.UUID `json:"case_id"`
	CaseName string    `json:"case_name"`
	Total    int       `json:"total_runs"`
	Passed   int       `json:"passed"`
	Failed   int       `json:"failed"`
}

type TrendPoint struct {
	Date    string `json:"date"`
	Passed  int    `json:"passed"`
	Failed  int    `json:"failed"`
	AvgDuration int64 `json:"avg_duration_ms"`
}

type DashboardService struct {
	repo DashboardRepository
}

func NewDashboardService(repo DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetSummary(ctx context.Context, projectID *uuid.UUID) (*DashboardSummary, error) {
	since := time.Now().AddDate(0, 0, -30) // Last 30 days
	return s.repo.GetSummary(ctx, projectID, since)
}

func (s *DashboardService) GetTrends(ctx context.Context, projectID *uuid.UUID, from, to time.Time) ([]TrendPoint, error) {
	return s.repo.GetTrends(ctx, projectID, from, to)
}

func (s *DashboardService) GetFlakyTests(ctx context.Context, projectID *uuid.UUID, limit int) ([]FlakyTest, error) {
	since := time.Now().AddDate(0, 0, -30) // Last 30 days
	return s.repo.GetFlakyTests(ctx, projectID, since, limit)
}
