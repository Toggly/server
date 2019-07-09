package app

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

// ErrorResponse creates {error: message} json body and responds with error code
func ErrorResponse(w http.ResponseWriter, r *http.Request, err error, code int) {
	render.Status(r, code)
	JSONResponse(w, r, map[string]interface{}{"error": err.Error()})
}

// JSONResponse creates json body
func JSONResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.JSON(w, r, data)
}

// NotFoundResponse creates empty json body and responds with 404 code
func NotFoundResponse(w http.ResponseWriter, r *http.Request, message string) {
	ErrorResponse(w, r, errors.New(message), http.StatusNotFound)
}

// UnauthorizedResponse creates empty json body and responds with 401 code
func UnauthorizedResponse(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusUnauthorized)
	render.PlainText(w, r, "")
}
