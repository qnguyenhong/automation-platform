package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/qnguyenhong/automation-platform/internal/model"
	"github.com/qnguyenhong/automation-platform/pkg/httputil"
)

func (d *Deps) ListNotificationConfigs(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	if projectIDStr == "" {
		httputil.Error(w, http.StatusBadRequest, "project_id is required")
		return
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid project ID")
		return
	}

	// This would call a notification service - for now return empty
	_ = projectID
	httputil.JSON(w, http.StatusOK, []model.NotificationConfig{})
}

func (d *Deps) CreateNotificationConfig(w http.ResponseWriter, r *http.Request) {
	var req model.CreateNotificationRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// TODO: implement notification service
	httputil.JSON(w, http.StatusCreated, req)
}

func (d *Deps) UpdateNotificationConfig(w http.ResponseWriter, r *http.Request) {
	configID, err := uuid.Parse(chi.URLParam(r, "configID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid config ID")
		return
	}

	var req model.UpdateNotificationRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// TODO: implement notification service
	_ = configID
	httputil.JSON(w, http.StatusOK, req)
}

func (d *Deps) DeleteNotificationConfig(w http.ResponseWriter, r *http.Request) {
	configID, err := uuid.Parse(chi.URLParam(r, "configID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid config ID")
		return
	}

	// TODO: implement notification service
	_ = configID
	w.WriteHeader(http.StatusNoContent)
}
