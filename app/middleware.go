package app

import (
	"context"
	"net/http"

	"bitbucket.org/toggly/toggly-server/models"
	"github.com/op/go-logging"
	"gopkg.in/toggly/go-utils.v2"
)

// CtxValue type
type CtxValue int

// CtxValue enum
const (
	CtxAPIVersion CtxValue = iota
	CtxValueOwner
	CtxValueRequestID
	CtxValueAuth
)

// Headers
const (
	XTogglyOwnerID string = "X-Toggly-Owner-Id"
)

// OwnerFromContext returns context value for project owner
func OwnerFromContext(r *http.Request) string {
	owner := r.Context().Value(CtxValueOwner)
	return owner.(string)
}

// OwnerCtx adds auth data to context
func OwnerCtx() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := GetLogger(r)
			owner := r.Header.Get(http.CanonicalHeaderKey(XTogglyOwnerID))
			if owner == "" {
				log.Error("Header X-Toggly-Owner-Id missed")
				NotFoundResponse(w, r, "Owner not found")
				return
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, CtxValueOwner, owner)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// GetLogger gets logger instance from context
func GetLogger(r *http.Request) *utils.StructuredLogger {
	log := r.Context().Value(models.ContextLoggerKey).(*logging.Logger)
	return &utils.StructuredLogger{Logger: log, R: r}
}
