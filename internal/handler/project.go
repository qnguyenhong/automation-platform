package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/qnguyenhong/automation-platform/internal/model"
	"github.com/qnguyenhong/automation-platform/internal/server/middleware"
	"github.com/qnguyenhong/automation-platform/pkg/httputil"
)

func (d *Deps) ListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := d.ProjectSvc.List(r.Context())
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, projects)
}

func (d *Deps) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req model.CreateProjectRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := middleware.GetUserID(r.Context())
	project, err := d.ProjectSvc.Create(r.Context(), req, userID)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, project)
}

func (d *Deps) GetProject(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid project ID")
		return
	}

	project, err := d.ProjectSvc.Get(r.Context(), id)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, project)
}

func (d *Deps) UpdateProject(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid project ID")
		return
	}

	var req model.UpdateProjectRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	project, err := d.ProjectSvc.Update(r.Context(), id, req)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, project)
}

func (d *Deps) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid project ID")
		return
	}

	if err := d.ProjectSvc.Delete(r.Context(), id); err != nil {
		httputil.AppError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
