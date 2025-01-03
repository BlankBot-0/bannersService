package handlers

import "banners/internal/usecase"

type Usecases struct {
	usecase.BannerManagementSystem // BMS interface
	usecase.AuthenticationSystem
	usecase.Cacher
}

// Controller - is controller/delivery layer
type Controller struct {
	Usecases
}

// NewController - returns Controller
func NewController(us Usecases) *Controller {
	return &Controller{
		Usecases: us,
	}
}
