package handler

import (
	"net/http"

	"github.com/qnguyenhong/automation-platform/internal/model"
	"github.com/qnguyenhong/automation-platform/internal/server/middleware"
	"github.com/qnguyenhong/automation-platform/pkg/httputil"
)

func (d *Deps) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := d.AuthService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, resp)
}

func (d *Deps) Register(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := d.AuthService.Register(r.Context(), req)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, user)
}

func (d *Deps) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	user, err := d.AuthService.GetUser(r.Context(), userID)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, user)
}
