package handlers

import (
	"banners/internal/usecase/authentification"
	"errors"
	"net/http"
	"strings"
)

func (c *Controller) AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("admin_token")
		if token == "" {
			ProcessError(w, ErrNoToken, http.StatusBadRequest)
			return
		}
		err := c.Usecases.AdminAuth(r.Context(), token)
		if errors.Is(err, authentification.ErrUnauthorized) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if errors.Is(err, authentification.ErrInvalidToken) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
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

func (c *Controller) UserAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		tokenSplit := strings.Fields(token)

		if len(tokenSplit) == 0 || tokenSplit[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if token == "" {
			ProcessError(w, ErrNoToken, http.StatusBadRequest)
			return
		}
		err := c.Usecases.UserAuth(r.Context(), token)
		if errors.Is(err, authentification.ErrInvalidToken) {
			ProcessError(w, err, http.StatusBadRequest)
			return
		}
		if errors.Is(err, authentification.ErrUnauthorized) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if errors.Is(err, authentification.ErrForbidden) {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r)
	})
}
