package app

import (
	"context"
	"net/http"

	"bitbucket.org/toggly/toggly-server/models"
	"bitbucket.org/toggly/toggly-server/service"
	"github.com/go-chi/chi"
	"github.com/op/go-logging"
	"gopkg.in/toggly/go-utils.v2"
)

// Headers
const (
	XTogglyOwnerID string = "X-Toggly-Owner-Id"
	XTogglyEnvID   string = "X-Toggly-Environment"
)

// OwnerCtx adds auth data to context
func OwnerCtx(defaultOwnerID string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := GetLogger(r)
			owner := r.Header.Get(http.CanonicalHeaderKey(XTogglyOwnerID))
			if owner == "" && defaultOwnerID != "" {
				owner = defaultOwnerID
			}
			if owner == "" {
				log.Error("Header X-Toggly-Owner-Id missed")
				models.NotFoundResponse(w, r, "Owner not found")
				return
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, models.CtxValueOwner, owner)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// EnvironmentCtx adds auth data to context
func EnvironmentCtx(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log := GetLogger(r)
		env := r.Header.Get(http.CanonicalHeaderKey(XTogglyEnvID))
		if env == "" {
			log.Error("Environemnt context is missed")
			models.ForbiddenResponse(w, r, "Unable to determine environment")
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, models.CtxValueEnvID, env)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// GetLogger gets logger instance from context
func GetLogger(r *http.Request) *utils.StructuredLogger {
	log := r.Context().Value(models.ContextLoggerKey).(*logging.Logger)
	return &utils.StructuredLogger{Logger: log, R: r}
}

// WithProjectCtx sets project id to context
func WithProjectCtx(srv *service.Project) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			code := chi.URLParam(r, "ProjectCode")
			srv.Logger = GetLogger(r)
			project, err := srv.Get(code)
			if err != nil {
				http.Error(w, http.StatusText(404), 404)
				return
			}
			ctx := context.WithValue(r.Context(), ContextProjectKey, project)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
