package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/qnguyenhong/automation-platform/internal/model"
	"github.com/qnguyenhong/automation-platform/internal/server/middleware"
	"github.com/qnguyenhong/automation-platform/pkg/httputil"
)

func (d *Deps) ListWorkers(w http.ResponseWriter, r *http.Request) {
	workers, err := d.WorkerSvc.List(r.Context())
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, workers)
}

func (d *Deps) GetWorker(w http.ResponseWriter, r *http.Request) {
	workerID, err := uuid.Parse(chi.URLParam(r, "workerID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid worker ID")
		return
	}

	worker, err := d.WorkerSvc.Get(r.Context(), workerID)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, worker)
}

func (d *Deps) DeleteWorker(w http.ResponseWriter, r *http.Request) {
	workerID, err := uuid.Parse(chi.URLParam(r, "workerID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid worker ID")
		return
	}

	if err := d.WorkerSvc.Delete(r.Context(), workerID); err != nil {
		httputil.AppError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (d *Deps) RegisterWorker(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterWorkerRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	worker, token, err := d.WorkerSvc.Register(r.Context(), req)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, map[string]any{
		"worker": worker,
		"token":  token,
	})
}

func (d *Deps) WorkerHeartbeat(w http.ResponseWriter, r *http.Request) {
	workerID, err := uuid.Parse(chi.URLParam(r, "workerID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid worker ID")
		return
	}

	var req model.HeartbeatRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := d.WorkerSvc.Heartbeat(r.Context(), workerID, req); err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (d *Deps) GetNextJob(w http.ResponseWriter, r *http.Request) {
	workerID := middleware.GetWorkerID(r.Context())
	workerUUID, _ := uuid.Parse(workerID)

	job, err := d.WorkerSvc.GetNextJob(r.Context(), workerUUID)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	if job == nil {
		httputil.JSON(w, http.StatusNoContent, nil)
		return
	}

	httputil.JSON(w, http.StatusOK, job)
}

func (d *Deps) SubmitJobResult(w http.ResponseWriter, r *http.Request) {
	jobID, err := uuid.Parse(chi.URLParam(r, "jobID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid job ID")
		return
	}

	var result model.JobResult
	if err := httputil.Decode(r, &result); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := d.RunSvc.SubmitResult(r.Context(), jobID, result); err != nil {
		httputil.AppError(w, err)
		return
	}

	// Notify via WebSocket
	d.Hub.Broadcast(map[string]any{
		"type":      "result_updated",
		"result_id": jobID.String(),
		"status":    result.Status,
	})

	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (d *Deps) SubmitJobLog(w http.ResponseWriter, r *http.Request) {
	jobID, err := uuid.Parse(chi.URLParam(r, "jobID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid job ID")
		return
	}

	var logEntry struct {
		Content string `json:"content"`
		Stream  string `json:"stream"`
	}
	if err := httputil.Decode(r, &logEntry); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	d.Hub.Broadcast(map[string]any{
		"type":    "log",
		"job_id":  jobID.String(),
		"content": logEntry.Content,
		"stream":  logEntry.Stream,
	})

	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
