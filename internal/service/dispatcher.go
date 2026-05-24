package service

import (
	"context"
	"log/slog"
	"sync"

	"github.com/qnguyenhong/automation-platform/internal/model"
)

// Dispatcher manages the job queue and dispatches jobs to workers.
type DispatcherImpl struct {
	mu       sync.Mutex
	queue    chan *model.WorkerJob
	logger   *slog.Logger
	workerSvc *WorkerService
}

func NewDispatcher(logger *slog.Logger, workerSvc *WorkerService, bufferSize int) *DispatcherImpl {
	if bufferSize < 1 {
		bufferSize = 100
	}
	return &DispatcherImpl{
		queue:     make(chan *model.WorkerJob, bufferSize),
		logger:    logger,
		workerSvc: workerSvc,
	}
}

// Enqueue adds a job to the dispatch queue.
func (d *DispatcherImpl) Enqueue(job *model.WorkerJob) {
	d.mu.Lock()
	defer d.mu.Unlock()

	select {
	case d.queue <- job:
		d.logger.Info("job enqueued", "job_id", job.JobID, "case", job.TestCase.Name)
	default:
		d.logger.Warn("job queue full, dropping job", "job_id", job.JobID)
	}
}

// Dequeue removes and returns the next job from the queue.
// Returns nil if the queue is empty.
func (d *DispatcherImpl) Dequeue() *model.WorkerJob {
	select {
	case job := <-d.queue:
		return job
	default:
		return nil
	}
}

// Start begins the dispatcher loop.
func (d *DispatcherImpl) Start(ctx context.Context) {
	d.logger.Info("dispatcher started")
	<-ctx.Done()
	d.logger.Info("dispatcher stopped")
}

// QueueSize returns the current number of jobs in the queue.
func (d *DispatcherImpl) QueueSize() int {
	return len(d.queue)
}
