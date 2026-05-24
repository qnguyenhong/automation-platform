package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/qnguyenhong/automation-platform/internal/model"
	"github.com/qnguyenhong/automation-platform/pkg/httputil"
)

func (d *Deps) ListTestCases(w http.ResponseWriter, r *http.Request) {
	suiteID, err := uuid.Parse(chi.URLParam(r, "suiteID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid suite ID")
		return
	}

	cases, err := d.TestCaseSvc.ListBySuite(r.Context(), suiteID)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, cases)
}

func (d *Deps) CreateTestCase(w http.ResponseWriter, r *http.Request) {
	suiteID, err := uuid.Parse(chi.URLParam(r, "suiteID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid suite ID")
		return
	}

	var req model.CreateTestCaseRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	testCase, err := d.TestCaseSvc.Create(r.Context(), suiteID, req)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, testCase)
}

func (d *Deps) GetTestCase(w http.ResponseWriter, r *http.Request) {
	caseID, err := uuid.Parse(chi.URLParam(r, "caseID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid case ID")
		return
	}

	testCase, err := d.TestCaseSvc.Get(r.Context(), caseID)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, testCase)
}

func (d *Deps) UpdateTestCase(w http.ResponseWriter, r *http.Request) {
	caseID, err := uuid.Parse(chi.URLParam(r, "caseID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid case ID")
		return
	}

	var req model.UpdateTestCaseRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	testCase, err := d.TestCaseSvc.Update(r.Context(), caseID, req)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, testCase)
}

func (d *Deps) DeleteTestCase(w http.ResponseWriter, r *http.Request) {
	caseID, err := uuid.Parse(chi.URLParam(r, "caseID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid case ID")
		return
	}

	if err := d.TestCaseSvc.Delete(r.Context(), caseID); err != nil {
		httputil.AppError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
