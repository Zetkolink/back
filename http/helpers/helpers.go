package helpers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-redis/redis/v8"
)

const (
	// APIPathSuffix is the path suffix for API endpoint URL.
	APIPathSuffix = "/api"
)

var (
	// APIVersionContextKey is context key for API version.
	APIVersionContextKey = &contextKey{"apiVersion"}

	// UserLoginKey is key for user login.
	UserLoginKey = &contextKey{"userLogin"}
)

type contextKey struct {
	name string
}

// ErrorResponse type represents error response.
type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Error      string `json:"error"`
}

// Render method is a rendering hook.
func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}

// NotFound method renders error with status code 404.
func NotFound(w http.ResponseWriter, r *http.Request, err error) {
	render.Render(w, r, NewErrorResponse(http.StatusNotFound, err))
}

// Conflict method renders error with status code 404.
func Conflict(w http.ResponseWriter, r *http.Request, err error) {
	render.Render(w, r, NewErrorResponse(http.StatusConflict, err))
}

// BadRequest method renders error with status code 400
func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	render.Render(w, r, NewErrorResponse(http.StatusBadRequest, err))
}

// Unauthorized method renders error with status code 401
func Unauthorized(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, NewErrorResponse(http.StatusUnauthorized,
		errors.New("401 Unauthorized")))
}

// Forbidden method renders error with status code 403
func Forbidden(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, NewErrorResponse(http.StatusForbidden,
		errors.New("403 Forbidden")))
}

// InternalServerError method renders error with status code 500.
func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err.Error())
	render.Render(w, r, NewErrorResponse(http.StatusInternalServerError,
		errors.New("500 Internal server error")))
}

// NewErrorResponse method creates new error response instance.
func NewErrorResponse(statusCode int, err error) *ErrorResponse {
	return &ErrorResponse{
		StatusCode: statusCode,
		Error:      err.Error(),
	}
}

// AccessController is a middleware for checking access privileges.
func AccessController(rdb *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handler := func(w http.ResponseWriter, r *http.Request) {
			var token string
			t := strings.Split(r.Header.Get("Authorization"), " ")

			if len(t) > 1 {
				token = t[1]
			}

			login, err := rdb.Get(r.Context(), token).Result()

			if err != nil {
				Unauthorized(w, r)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, UserLoginKey, login)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(handler)
	}
}
