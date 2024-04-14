package handlers

import (
	"banners/internal/config"
	"net/http"
)

func (c *Controller) NewServer(cfg config.HTTPServer) http.Server {
	// Router layer
	router := http.NewServeMux()

	userRouter := http.NewServeMux()
	userRouter.HandleFunc("GET /user_banner", c.UserBannerHandler)

	adminRouter := http.NewServeMux()
	adminRouter.HandleFunc("GET /banner", c.BannersSortedHandler)
	adminRouter.HandleFunc("POST /banner", c.CreateBannerHandler)
	adminRouter.HandleFunc("PATCH /banner/{id}", c.PatchBannerHandler)
	adminRouter.HandleFunc("DELETE /banner/{id}", c.DeleteBannerHandler)
	adminRouter.HandleFunc("GET /banner/versions", c.BannerVersionsHandler)

	router.Handle("GET /banner", c.AdminAuthMiddleware(adminRouter))
	router.Handle("POST /banner", c.AdminAuthMiddleware(adminRouter))
	router.Handle("PATCH /banner/{id}", c.AdminAuthMiddleware(adminRouter))
	router.Handle("DELETE /banner/{id}", c.AdminAuthMiddleware(adminRouter))
	router.Handle("GET /banner/versions", c.AdminAuthMiddleware(adminRouter))

	router.HandleFunc("GET /user_token", c.UserToken)
	router.HandleFunc("GET /admin_token", c.AdminToken)

	router.Handle("/", c.UserAuthMiddleware(userRouter))

	return http.Server{
		Addr:    cfg.Address,
		Handler: Logging(router),
	}
}
