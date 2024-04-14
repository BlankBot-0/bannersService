package handlers

import (
	"banners/internal/usecase/authentification"
	"errors"
	"net/http"
)

func AdminAuth(next http.Handler, c *Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("admin_token")
		if token == "" {
			http.Error(w, "Missing token", http.StatusBadRequest)
		}
		err := c.Usecases.AdminAuth(r.Context(), token)
		if errors.Is(err, authentification.ErrUnauthorized) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if errors.Is(err, authentification.ErrInvalidToken) {
			ProcessError(w, err, http.StatusBadRequest)
			return
		}
		if errors.Is(err, authentification.ErrForbidden) {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if err != nil {
			ProcessError(w, ErrInternal, http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func UserAuth(next http.Handler, c *Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("user_token")
		if token == "" {
			http.Error(w, "Missing token", http.StatusBadRequest)
		}
		err := c.Usecases.UserAuth(r.Context(), token)
		if errors.Is(err, authentification.ErrInvalidToken) {
			ProcessError(w, err, http.StatusBadRequest)
		}
		if errors.Is(err, authentification.ErrUnauthorized) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		if errors.Is(err, authentification.ErrForbidden) {
			http.Error(w, err.Error(), http.StatusForbidden)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r)
	})
}
