package handlers

import (
	"banners/internal/config"
	"net/http"
)

func (c *Controller) NewServer(cfg config.HTTPServer) http.Server {
	// Router layer
	router := http.NewServeMux()

	userRouter := http.NewServeMux()
	userRouter.HandleFunc("GET /user_token", c.UserToken)
	userRouter.HandleFunc("GET /user_banner", c.UserBannerHandler)

	adminRouter := http.NewServeMux()
	adminRouter.HandleFunc("GET /admin_token", c.AdminToken)
	adminRouter.HandleFunc("GET /banner", c.BannersSortedHandler)
	adminRouter.HandleFunc("POST /banner", c.CreateBannerHandler)
	adminRouter.HandleFunc("PATCH /banner/{id}", c.PatchBannerHandler)
	adminRouter.HandleFunc("DELETE /banner/{id}", c.DeleteBannerHandler)
	adminRouter.HandleFunc("GET /banner/versions", c.BannerVersionsHandler)

	router.Handle("GET /banner", AdminAuth(adminRouter, c))
	router.Handle("POST /banner", AdminAuth(adminRouter, c))
	router.Handle("PATCH /banner/{id}", AdminAuth(adminRouter, c))
	router.Handle("DELETE /banner/{id}", AdminAuth(adminRouter, c))
	router.Handle("GET /banner/versions", AdminAuth(adminRouter, c))

	router.Handle("GET /user_banner", UserAuth(userRouter, c))

	return http.Server{
		Addr:    cfg.Address,
		Handler: Logging(router),
	}
}
