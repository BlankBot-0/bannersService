package handlers

import (
	"net/http"
)

// NewRouter - returns http.Handler
func (c *Controller) NewRouter() http.Handler {
	// Router layer
	mux := http.NewServeMux()

	mux.HandleFunc("GET /user_banner", c.UserBannerHandler)
	mux.HandleFunc("GET /banner", c.BannersSortedHandler)
	mux.HandleFunc("POST /banner", c.CreateBannerHandler)
	mux.HandleFunc("PATCH /banner/{id}", c.PatchBannerHandler)
	mux.HandleFunc("DELETE /banner/{id}", c.DeleteBannerHandler)

	return mux
}
