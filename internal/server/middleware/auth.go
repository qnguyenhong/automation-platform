package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/qnguyenhong/automation-platform/internal/service"
	"github.com/qnguyenhong/automation-platform/pkg/httputil"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
	WorkerIDKey contextKey = "worker_id"
)

// Auth middleware validates JWT token or API key.
func Auth(authSvc *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var token string

			// Check Authorization header (Bearer token)
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			}

			// Check X-API-Key header
			if token == "" {
				token = r.Header.Get("X-API-Key")
			}

			// Check query parameter "token" (especially for WebSockets)
			if token == "" {
				token = r.URL.Query().Get("token")
			}

			if token == "" {
				httputil.Error(w, http.StatusUnauthorized, "missing authentication")
				return
			}

			userID, role, err := authSvc.ValidateToken(r.Context(), token)
			if err != nil {
				httputil.Error(w, http.StatusUnauthorized, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UserRoleKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// WorkerAuth middleware validates worker registration token.
func WorkerAuth(workerSvc *service.WorkerService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Worker registration endpoint doesn't require auth
			if r.URL.Path == "/api/v1/workers/register" && r.Method == "POST" {
				next.ServeHTTP(w, r)
				return
			}

			token := r.Header.Get("X-Worker-Token")
			if token == "" {
				httputil.Error(w, http.StatusUnauthorized, "missing worker token")
				return
			}

			workerID, err := workerSvc.ValidateWorkerToken(r.Context(), token)
			if err != nil {
				httputil.Error(w, http.StatusUnauthorized, "invalid worker token")
				return
			}

			ctx := context.WithValue(r.Context(), WorkerIDKey, workerID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts user ID from context.
func GetUserID(ctx context.Context) string {
	if v, ok := ctx.Value(UserIDKey).(string); ok {
		return v
	}
	return ""
}

// GetUserRole extracts user role from context.
func GetUserRole(ctx context.Context) string {
	if v, ok := ctx.Value(UserRoleKey).(string); ok {
		return v
	}
	return ""
}

// GetWorkerID extracts worker ID from context.
func GetWorkerID(ctx context.Context) string {
	if v, ok := ctx.Value(WorkerIDKey).(string); ok {
		return v
	}
	return ""
}
