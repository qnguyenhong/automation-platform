package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/qnguyenhong/automation-platform/internal/model"
	"github.com/qnguyenhong/automation-platform/internal/service"
	"github.com/qnguyenhong/automation-platform/pkg/httputil"
)

// --- Environment Handlers ---

func (d *Deps) ListEnvironments(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid project ID")
		return
	}
	envs, err := d.EnvSvc.List(r.Context(), projectID)
	if err != nil {
		httputil.AppError(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, envs)
}

func (d *Deps) CreateEnvironment(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid project ID")
		return
	}
	var req model.CreateEnvironmentRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	env, err := d.EnvSvc.Create(r.Context(), projectID, req)
	if err != nil {
		httputil.AppError(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, env)
}

func (d *Deps) GetEnvironment(w http.ResponseWriter, r *http.Request) {
	envID, err := uuid.Parse(chi.URLParam(r, "envID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid environment ID")
		return
	}
	env, err := d.EnvSvc.Get(r.Context(), envID)
	if err != nil {
		httputil.AppError(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, env)
}

func (d *Deps) UpdateEnvironment(w http.ResponseWriter, r *http.Request) {
	envID, err := uuid.Parse(chi.URLParam(r, "envID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid environment ID")
		return
	}
	var req model.UpdateEnvironmentRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	env, err := d.EnvSvc.Update(r.Context(), envID, req)
	if err != nil {
		httputil.AppError(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, env)
}

func (d *Deps) DeleteEnvironment(w http.ResponseWriter, r *http.Request) {
	envID, err := uuid.Parse(chi.URLParam(r, "envID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid environment ID")
		return
	}
	if err := d.EnvSvc.Delete(r.Context(), envID); err != nil {
		httputil.AppError(w, err)
		return
	}
	httputil.JSON(w, http.StatusNoContent, nil)
}

// --- Dataset Handlers ---

func (d *Deps) ListDatasets(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid project ID")
		return
	}
	datasets, err := d.DatasetSvc.List(r.Context(), projectID)
	if err != nil {
		httputil.AppError(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, datasets)
}

func (d *Deps) CreateDataset(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid project ID")
		return
	}
	var req model.CreateDatasetRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	ds, err := d.DatasetSvc.Create(r.Context(), projectID, req)
	if err != nil {
		httputil.AppError(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, ds)
}

func (d *Deps) GetDataset(w http.ResponseWriter, r *http.Request) {
	dsID, err := uuid.Parse(chi.URLParam(r, "datasetID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid dataset ID")
		return
	}
	ds, err := d.DatasetSvc.Get(r.Context(), dsID)
	if err != nil {
		httputil.AppError(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, ds)
}

func (d *Deps) UpdateDataset(w http.ResponseWriter, r *http.Request) {
	dsID, err := uuid.Parse(chi.URLParam(r, "datasetID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid dataset ID")
		return
	}
	var req model.UpdateDatasetRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	ds, err := d.DatasetSvc.Update(r.Context(), dsID, req)
	if err != nil {
		httputil.AppError(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, ds)
}

func (d *Deps) DeleteDataset(w http.ResponseWriter, r *http.Request) {
	dsID, err := uuid.Parse(chi.URLParam(r, "datasetID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid dataset ID")
		return
	}
	if err := d.DatasetSvc.Delete(r.Context(), dsID); err != nil {
		httputil.AppError(w, err)
		return
	}
	httputil.JSON(w, http.StatusNoContent, nil)
}

// --- Webhook Handlers ---

func (d *Deps) ListWebhooks(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid project ID")
		return
	}
	webhooks, err := d.WebhookSvc.List(r.Context(), projectID)
	if err != nil {
		httputil.AppError(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, webhooks)
}

func (d *Deps) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid project ID")
		return
	}
	var req model.CreateWebhookRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	wh, err := d.WebhookSvc.Create(r.Context(), projectID, req)
	if err != nil {
		httputil.AppError(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, wh)
}

func (d *Deps) GetWebhook(w http.ResponseWriter, r *http.Request) {
	whID, err := uuid.Parse(chi.URLParam(r, "webhookID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid webhook ID")
		return
	}
	wh, err := d.WebhookSvc.Get(r.Context(), whID)
	if err != nil {
		httputil.AppError(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, wh)
}

func (d *Deps) UpdateWebhook(w http.ResponseWriter, r *http.Request) {
	whID, err := uuid.Parse(chi.URLParam(r, "webhookID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid webhook ID")
		return
	}
	var req model.UpdateWebhookRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	wh, err := d.WebhookSvc.Update(r.Context(), whID, req)
	if err != nil {
		httputil.AppError(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, wh)
}

func (d *Deps) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	whID, err := uuid.Parse(chi.URLParam(r, "webhookID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid webhook ID")
		return
	}
	if err := d.WebhookSvc.Delete(r.Context(), whID); err != nil {
		httputil.AppError(w, err)
		return
	}
	httputil.JSON(w, http.StatusNoContent, nil)
}

// TriggerWebhook is a public endpoint that triggers a run via webhook secret.
func (d *Deps) TriggerWebhook(w http.ResponseWriter, r *http.Request) {
	secret := r.Header.Get("X-Webhook-Secret")
	if secret == "" {
		// Try query param
		secret = r.URL.Query().Get("secret")
	}
	if secret == "" {
		httputil.Error(w, http.StatusUnauthorized, "missing webhook secret")
		return
	}

	run, err := d.WebhookSvc.Trigger(r.Context(), secret)
	if err != nil {
		httputil.AppError(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, run)
}

// --- Export Handler ---

func (d *Deps) ExportSuite(w http.ResponseWriter, r *http.Request) {
	var req service.ExportRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	result, err := d.ExportSvc.Export(r.Context(), req)
	if err != nil {
		httputil.AppError(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, result)
}
