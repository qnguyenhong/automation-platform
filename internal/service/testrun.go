package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/qnguyenhong/automation-platform/internal/model"
	"github.com/qnguyenhong/automation-platform/pkg/errors"
)

type TestRunRepository interface {
	Create(ctx context.Context, suiteID, projectID uuid.UUID, status, trigger string, triggeredBy *uuid.UUID, totalCases int, metadata []byte) (*model.TestRun, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.TestRun, error)
	ListAll(ctx context.Context, limit, offset int) ([]model.TestRun, int64, error)
	ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]model.TestRun, int64, error)
	ListBySuite(ctx context.Context, suiteID uuid.UUID, limit, offset int) ([]model.TestRun, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, startedAt, finishedAt *time.Time, durationMS *int64) (*model.TestRun, error)
	UpdateCounts(ctx context.Context, id uuid.UUID, passed, failed, skipped, errored int) (*model.TestRun, error)
}

type TestResultRepository interface {
	Create(ctx context.Context, runID, caseID uuid.UUID, workerID *uuid.UUID, status string) (*model.TestResult, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.TestResult, error)
	ListByRun(ctx context.Context, runID uuid.UUID) ([]model.TestResult, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, startedAt, finishedAt *time.Time, durationMS *int64) (*model.TestResult, error)
	UpdateData(ctx context.Context, id uuid.UUID, status string, errorMessage *string, assertions, requestData, responseData, artifacts, metrics []byte, stdout, stderr *string, durationMS *int64, finishedAt *time.Time) (*model.TestResult, error)
}

type LoadMetricsRepository interface {
	Create(ctx context.Context, metrics *model.LoadMetricsDb) error
	GetByRun(ctx context.Context, runID uuid.UUID) (*model.LoadMetricsDb, error)
	GetByResult(ctx context.Context, resultID uuid.UUID) (*model.LoadMetricsDb, error)
}

type Dispatcher interface {
	Enqueue(job *model.WorkerJob)
}

type TestRunService struct {
	runRepo    TestRunRepository
	resultRepo TestResultRepository
	suiteRepo  TestSuiteRepository
	caseRepo   TestCaseRepository
	loadRepo   LoadMetricsRepository
	dispatcher Dispatcher
}

func NewTestRunService(runRepo TestRunRepository, resultRepo TestResultRepository, suiteRepo TestSuiteRepository, caseRepo TestCaseRepository, loadRepo LoadMetricsRepository, dispatcher Dispatcher) *TestRunService {
	return &TestRunService{
		runRepo:    runRepo,
		resultRepo: resultRepo,
		suiteRepo:  suiteRepo,
		caseRepo:   caseRepo,
		loadRepo:   loadRepo,
		dispatcher: dispatcher,
	}
}

func (s *TestRunService) Trigger(ctx context.Context, suiteID uuid.UUID, userID string, req model.TriggerRunRequest) (*model.TestRun, error) {
	suite, err := s.suiteRepo.GetByID(ctx, suiteID)
	if err != nil {
		return nil, err
	}

	cases, err := s.caseRepo.ListBySuite(ctx, suiteID)
	if err != nil {
		return nil, err
	}

	enabledCases := 0
	for _, c := range cases {
		if c.Enabled {
			enabledCases++
		}
	}

	if enabledCases == 0 {
		return nil, errors.NewAppError("BAD_REQUEST", "no enabled test cases in suite", nil)
	}

	trigger := model.TriggerManual
	if req.Trigger != "" {
		trigger = req.Trigger
	}

	var triggeredBy *uuid.UUID
	if userID != "" {
		uid, _ := uuid.Parse(userID)
		if uid != uuid.Nil {
			triggeredBy = &uid
		}
	}

	metadata := req.Metadata
	if metadata == nil {
		metadata = json.RawMessage(`{}`)
	}

	run, err := s.runRepo.Create(ctx, suiteID, suite.ProjectID, string(model.StatusPending), string(trigger), triggeredBy, enabledCases, metadata)
	if err != nil {
		return nil, err
	}

	// Create test results for each enabled case
	for _, c := range cases {
		if !c.Enabled {
			continue
		}
		result, err := s.resultRepo.Create(ctx, run.ID, c.ID, nil, string(model.StatusPending))
		if err != nil {
			return nil, err
		}

		// Enqueue job for dispatcher
		s.dispatcher.Enqueue(&model.WorkerJob{
			JobID:    result.ID.String(),
			ResultID: result.ID.String(),
			TestCase: c,
			Suite:    *suite,
			Run:      *run,
			Timeout:  5 * time.Minute,
		})
	}

	// Update run status to running
	now := time.Now()
	run, err = s.runRepo.UpdateStatus(ctx, run.ID, string(model.StatusRunning), &now, nil, nil)
	if err != nil {
		return nil, err
	}

	return run, nil
}

func (s *TestRunService) Get(ctx context.Context, id uuid.UUID) (*model.TestRun, error) {
	return s.runRepo.GetByID(ctx, id)
}

func (s *TestRunService) ListByProject(ctx context.Context, projectID uuid.UUID, page, perPage int) ([]model.TestRun, int64, error) {
	offset := (page - 1) * perPage
	return s.runRepo.ListByProject(ctx, projectID, perPage, offset)
}

func (s *TestRunService) ListAll(ctx context.Context, page, perPage int) ([]model.TestRun, int64, error) {
	offset := (page - 1) * perPage
	return s.runRepo.ListAll(ctx, perPage, offset)
}

func (s *TestRunService) Cancel(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	_, err := s.runRepo.UpdateStatus(ctx, id, string(model.StatusCancelled), nil, &now, nil)
	return err
}

func (s *TestRunService) ListResults(ctx context.Context, runID uuid.UUID) ([]model.TestResult, error) {
	return s.resultRepo.ListByRun(ctx, runID)
}

func (s *TestRunService) GetResult(ctx context.Context, id uuid.UUID) (*model.TestResult, error) {
	return s.resultRepo.GetByID(ctx, id)
}

func (s *TestRunService) SubmitResult(ctx context.Context, resultID uuid.UUID, result model.JobResult) error {
	now := time.Now()

	errorMsg := &result.ErrorMessage
	if result.ErrorMessage == "" {
		errorMsg = nil
	}

	_, err := s.resultRepo.UpdateData(
		ctx,
		resultID,
		string(result.Status),
		errorMsg,
		result.Assertions,
		result.RequestData,
		result.ResponseData,
		result.Artifacts,
		result.Metrics,
		&result.Stdout,
		&result.Stderr,
		&result.DurationMS,
		&now,
	)
	if err != nil {
		return err
	}

	// Get the result to find the run ID and case
	testResult, err := s.resultRepo.GetByID(ctx, resultID)
	if err != nil {
		return err
	}

	// Persist load metrics if test suite type is load
	testCase, err := s.caseRepo.GetByID(ctx, testResult.CaseID)
	if err == nil {
		suite, err := s.suiteRepo.GetByID(ctx, testCase.SuiteID)
		if err == nil && suite.TestType == "load" && len(result.Metrics) > 0 {
			var rawMetrics model.LoadMetrics
			if err := json.Unmarshal(result.Metrics, &rawMetrics); err == nil {
				type TimeSeriesPayload struct {
					RPS []model.TimePoint `json:"rps"`
					P95 []model.TimePoint `json:"p95"`
				}
				tsPayload := TimeSeriesPayload{
					RPS: rawMetrics.TimeSeriesRPS,
					P95: rawMetrics.TimeSeriesP95,
				}
				tsBytes, _ := json.Marshal(tsPayload)
				statusCodesBytes, _ := json.Marshal(rawMetrics.StatusCodes)

				dbMetrics := &model.LoadMetricsDb{
					ResultID:      testResult.ID,
					RunID:         testResult.RunID,
					TotalRequests: rawMetrics.TotalRequests,
					SuccessCount:  rawMetrics.SuccessCount,
					ErrorCount:    rawMetrics.ErrorCount,
					ErrorRate:     rawMetrics.ErrorRate,
					ThroughputRPS: rawMetrics.ThroughputRPS,
					MinLatencyMS:  rawMetrics.MinLatencyMS,
					MaxLatencyMS:  rawMetrics.MaxLatencyMS,
					AvgLatencyMS:  rawMetrics.AvgLatencyMS,
					P50LatencyMS:  rawMetrics.P50LatencyMS,
					P90LatencyMS:  rawMetrics.P90LatencyMS,
					P95LatencyMS:  rawMetrics.P95LatencyMS,
					P99LatencyMS:  rawMetrics.P99LatencyMS,
					TotalBytes:    rawMetrics.TotalBytes,
					TimeSeries:    json.RawMessage(tsBytes),
					StatusCodes:   json.RawMessage(statusCodesBytes),
				}
				_ = s.loadRepo.Create(ctx, dbMetrics)
			}
		}
	}

	// Check if we should retry (only for failed/error results)
	if result.Status == model.StatusFailed || result.Status == model.StatusError {
		if retried := s.attemptRetry(ctx, testResult); retried {
			return nil // Don't update run counts yet, wait for retry
		}
	}

	// Update run counts
	counts, err := s.resultRepo.ListByRun(ctx, testResult.RunID)
	if err != nil {
		return err
	}

	var passed, failed, skipped, errored int
	allDone := true
	for _, r := range counts {
		switch model.StatusType(r.Status) {
		case model.StatusPassed:
			passed++
		case model.StatusFailed:
			failed++
		case model.StatusSkipped:
			skipped++
		case model.StatusError:
			errored++
		default:
			allDone = false
		}
	}

	_, err = s.runRepo.UpdateCounts(ctx, testResult.RunID, passed, failed, skipped, errored)
	if err != nil {
		return err
	}

	// If all results are done, mark run as complete
	if allDone {
		status := model.StatusPassed
		if failed > 0 || errored > 0 {
			status = model.StatusFailed
		}

		run, err := s.runRepo.GetByID(ctx, testResult.RunID)
		if err != nil {
			return err
		}

		duration := now.Sub(*run.StartedAt).Milliseconds()
		_, err = s.runRepo.UpdateStatus(ctx, testResult.RunID, string(status), nil, &now, &duration)
		if err != nil {
			return err
		}
	}

	return nil
}

// attemptRetry checks if a failed result should be retried and enqueues a retry job if so.
func (s *TestRunService) attemptRetry(ctx context.Context, testResult *model.TestResult) bool {
	// Get the test case to check retry config
	testCase, err := s.caseRepo.GetByID(ctx, testResult.CaseID)
	if err != nil {
		return false
	}

	// Get suite for retry config fallback
	suite, err := s.suiteRepo.GetByID(ctx, testCase.SuiteID)
	if err != nil {
		return false
	}

	// Determine max retries (case-level overrides suite-level)
	// This would require the new columns; for now check config JSON
	var caseConfig map[string]any
	if err := json.Unmarshal(testCase.Config, &caseConfig); err == nil {
		if maxRetries, ok := caseConfig["max_retries"]; ok {
			if mr, ok := maxRetries.(float64); ok && int(mr) > 0 {
				// Count existing retries for this case in this run
				results, _ := s.resultRepo.ListByRun(ctx, testResult.RunID)
				retryCount := 0
				for _, r := range results {
					if r.CaseID == testResult.CaseID && r.RetryOf != nil {
						retryCount++
					}
				}

				if retryCount < int(mr) {
					// Create a new retry result
					newResult, err := s.resultRepo.Create(ctx, testResult.RunID, testResult.CaseID, nil, string(model.StatusPending))
					if err != nil {
						return false
					}

					// Enqueue retry job
					run, _ := s.runRepo.GetByID(ctx, testResult.RunID)
					if run != nil {
						s.dispatcher.Enqueue(&model.WorkerJob{
							JobID:    newResult.ID.String(),
							ResultID: newResult.ID.String(),
							TestCase: *testCase,
							Suite:    *suite,
							Run:      *run,
							Timeout:  5 * time.Minute,
						})
					}
					return true
				}
			}
		}
	}

	_ = suite // suppress unused
	return false
}
