package handlers

import (
	"banners/internal/config"
	"net/http"
)

func (c *Controller) NewServer(cfg config.HTTPServer) http.Server {
	// Router layer
	router := http.NewServeMux()

	adminRouter := http.NewServeMux()
	adminRouter.HandleFunc("GET /banner", c.BannersSortedHandler)
	adminRouter.HandleFunc("POST /banner", c.CreateBannerHandler)
	adminRouter.HandleFunc("PATCH /banner/{id}", c.PatchBannerHandler)
	adminRouter.HandleFunc("DELETE /banner/{id}", c.DeleteBannerHandler)
	adminRouter.HandleFunc("GET /banner/versions", c.BannerVersionsHandler)

	router.Handle("/", c.AdminAuthMiddleware(adminRouter))

	router.HandleFunc("GET /user_token", c.UserToken)
	router.HandleFunc("GET /admin_token", c.AdminToken)

	router.Handle("GET /user_banner", c.UserAuthMiddleware(http.HandlerFunc(c.UserBannerHandler)))

	return http.Server{
		Addr:    cfg.Address,
		Handler: Logging(router),
	}
}
