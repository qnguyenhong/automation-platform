package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/qnguyenhong/automation-platform/pkg/httputil"
)

func (d *Deps) GetRunLoadMetrics(w http.ResponseWriter, r *http.Request) {
	runID, err := uuid.Parse(chi.URLParam(r, "runID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid run ID")
		return
	}

	metrics, err := d.LoadRepo.GetByRun(r.Context(), runID)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, metrics)
}

func (d *Deps) GetResultLoadMetrics(w http.ResponseWriter, r *http.Request) {
	resultID, err := uuid.Parse(chi.URLParam(r, "resultID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid result ID")
		return
	}

	metrics, err := d.LoadRepo.GetByResult(r.Context(), resultID)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, metrics)
}
