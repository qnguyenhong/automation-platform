package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/qnguyenhong/automation-platform/pkg/httputil"
)

func (d *Deps) GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	var projectID *uuid.UUID
	if projectIDStr != "" {
		id, err := uuid.Parse(projectIDStr)
		if err == nil {
			projectID = &id
		}
	}

	summary, err := d.DashboardSvc.GetSummary(r.Context(), projectID)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, summary)
}

func (d *Deps) GetDashboardTrends(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	var projectID *uuid.UUID
	if projectIDStr != "" {
		id, err := uuid.Parse(projectIDStr)
		if err == nil {
			projectID = &id
		}
	}

	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	if days < 1 {
		days = 30
	}

	from := time.Now().AddDate(0, 0, -days)
	to := time.Now()

	trends, err := d.DashboardSvc.GetTrends(r.Context(), projectID, from, to)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, trends)
}

func (d *Deps) GetFlakyTests(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	var projectID *uuid.UUID
	if projectIDStr != "" {
		id, err := uuid.Parse(projectIDStr)
		if err == nil {
			projectID = &id
		}
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 20
	}

	flakyTests, err := d.DashboardSvc.GetFlakyTests(r.Context(), projectID, limit)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, flakyTests)
}
