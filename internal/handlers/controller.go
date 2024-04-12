package handlers

import "banners/internal/usecase"

type Usecases struct {
	usecase.BannerManagementSystem // OMS interface
}

// Controller - is controller/delivery layer
type Controller struct {
	Usecases
	/* ... */
}

// NewController - returns Controller
func NewController(us Usecases) *Controller {
	return &Controller{
		Usecases: us,
	}
}
