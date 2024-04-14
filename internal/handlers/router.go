package handlers

import (
	"banners/internal/config"
	"banners/internal/middleware"
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

	router.Handle("/", middleware.EnsureAdmin(adminRouter))
	router.Handle("GET /user_banner", middleware.UserAuth(userRouter))

	return http.Server{
		Addr:    cfg.Address,
		Handler: middleware.Logging(router),
	}
}
