package models

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/op/go-logging"
)

// ErrorResponseWithStatus creates {error: message} json body and responds with error code
func ErrorResponseWithStatus(w http.ResponseWriter, r *http.Request, err error, code int) {
	render.Status(r, code)
	JSONResponse(w, r, map[string]interface{}{"error": err.Error()})
}

// ErrorResponse creates {error: message} json body and responds with error code
func ErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	if err, ok := err.(*ErrStatusedResponse); ok {
		render.Status(r, err.Code)
		JSONResponse(w, r, map[string]interface{}{"error": err.Error()})
	}
}

// JSONResponse creates json body
func JSONResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.JSON(w, r, data)
}

// NotFoundResponse creates empty json body and responds with 404 code
func NotFoundResponse(w http.ResponseWriter, r *http.Request, message string) {
	log := r.Context().Value(ContextLoggerKey).(*logging.Logger)
	log.Error(message)
	ErrorResponseWithStatus(w, r, errors.New(message), http.StatusNotFound)
}

// ForbiddenResponse creates empty json body and responds with 403 code
func ForbiddenResponse(w http.ResponseWriter, r *http.Request, message string) {
	ErrorResponseWithStatus(w, r, errors.New(message), http.StatusForbidden)
}

// UnauthorizedResponse creates empty json body and responds with 401 code
func UnauthorizedResponse(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusUnauthorized)
	render.PlainText(w, r, "")
}

// ErrStatusedResponse type
type ErrStatusedResponse struct {
	Code    int
	Message string
}

// ErrBadRequest func
func ErrBadRequest(message string) *ErrStatusedResponse {
	return &ErrStatusedResponse{Message: message, Code: http.StatusBadRequest}
}

// ErrInternalServer func
func ErrInternalServer(message string) *ErrStatusedResponse {
	return &ErrStatusedResponse{Message: message, Code: http.StatusInternalServerError}
}

// ErrNotFound func
func ErrNotFound(message string) *ErrStatusedResponse {
	return &ErrStatusedResponse{Message: message, Code: http.StatusNotFound}
}

// ErrConflict func
func ErrConflict(message string) *ErrStatusedResponse {
	return &ErrStatusedResponse{Message: message, Code: http.StatusConflict}
}

func (e *ErrStatusedResponse) Error() string {
	return e.Message
}
