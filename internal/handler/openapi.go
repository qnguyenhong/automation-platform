package handler

import (
	"context"
	"net/http"

	"github.com/qnguyenhong/automation-platform/internal/server/middleware"
	"github.com/qnguyenhong/automation-platform/internal/service"
	"github.com/qnguyenhong/automation-platform/pkg/httputil"
)

func getCtxUserID(ctx context.Context) string {
	return middleware.GetUserID(ctx)
}

func (d *Deps) ParseOpenAPI(w http.ResponseWriter, r *http.Request) {
	var req service.ParseOpenAPIRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := d.OpenAPISvc.ParseSpec(r.Context(), req.Content)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, result)
}

func (d *Deps) ImportOpenAPI(w http.ResponseWriter, r *http.Request) {
	var req service.ImportOpenAPIRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := getCtxUserID(r.Context())
	result, err := d.OpenAPISvc.ImportEndpoints(r.Context(), req, userID)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, result)
}
