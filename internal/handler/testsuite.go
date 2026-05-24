package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/qnguyenhong/automation-platform/internal/model"
	"github.com/qnguyenhong/automation-platform/pkg/httputil"
)

func (d *Deps) ListTestSuites(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid project ID")
		return
	}

	suites, err := d.SuiteSvc.ListByProject(r.Context(), projectID)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, suites)
}

func (d *Deps) CreateTestSuite(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid project ID")
		return
	}

	var req model.CreateTestSuiteRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	suite, err := d.SuiteSvc.Create(r.Context(), projectID, req)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, suite)
}

func (d *Deps) GetTestSuite(w http.ResponseWriter, r *http.Request) {
	suiteID, err := uuid.Parse(chi.URLParam(r, "suiteID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid suite ID")
		return
	}

	suite, err := d.SuiteSvc.Get(r.Context(), suiteID)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, suite)
}

func (d *Deps) UpdateTestSuite(w http.ResponseWriter, r *http.Request) {
	suiteID, err := uuid.Parse(chi.URLParam(r, "suiteID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid suite ID")
		return
	}

	var req model.UpdateTestSuiteRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	suite, err := d.SuiteSvc.Update(r.Context(), suiteID, req)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, suite)
}

func (d *Deps) DeleteTestSuite(w http.ResponseWriter, r *http.Request) {
	suiteID, err := uuid.Parse(chi.URLParam(r, "suiteID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid suite ID")
		return
	}

	if err := d.SuiteSvc.Delete(r.Context(), suiteID); err != nil {
		httputil.AppError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
