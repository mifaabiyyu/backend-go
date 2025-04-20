package utils

import (
	"net/http"

	"go.uber.org/zap"
)

type AppWrapper struct {
	Logger *zap.SugaredLogger
}

func (app *AppWrapper) InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.Logger.Errorw("internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *AppWrapper) ForbiddenResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.Logger.Warnw("forbidden", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteJSONError(w, http.StatusForbidden, err.Error())
}

func (app *AppWrapper) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.Logger.Warnf("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *AppWrapper) ConflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.Logger.Errorf("conflict response", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteJSONError(w, http.StatusConflict, err.Error())
}

func (app *AppWrapper) NotFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.Logger.Warnf("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteJSONError(w, http.StatusNotFound, "not found")
}

func (app *AppWrapper) UnauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.Logger.Warnf("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteJSONError(w, http.StatusUnauthorized, `unauthorized : `+err.Error())
}

func (app *AppWrapper) UnauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.Logger.Warnf("unauthorized basic error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *AppWrapper) RateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.Logger.Warnw("rate limit exceeded", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("Retry-After", retryAfter)

	WriteJSONError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}
