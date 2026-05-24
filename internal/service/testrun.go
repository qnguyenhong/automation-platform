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

type Dispatcher interface {
	Enqueue(job *model.WorkerJob)
}

type TestRunService struct {
	runRepo    TestRunRepository
	resultRepo TestResultRepository
	suiteRepo  TestSuiteRepository
	caseRepo   TestCaseRepository
	dispatcher Dispatcher
}

func NewTestRunService(runRepo TestRunRepository, resultRepo TestResultRepository, suiteRepo TestSuiteRepository, caseRepo TestCaseRepository, dispatcher Dispatcher) *TestRunService {
	return &TestRunService{
		runRepo:    runRepo,
		resultRepo: resultRepo,
		suiteRepo:  suiteRepo,
		caseRepo:   caseRepo,
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

	// Get the result to find the run ID
	testResult, err := s.resultRepo.GetByID(ctx, resultID)
	if err != nil {
		return err
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
