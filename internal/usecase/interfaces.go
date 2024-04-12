package usecase

import (
	"banners/internal/models"
	"context"
)

type BannerManagementSystem interface {
	CreateBanner(ctx context.Context, banner models.Banner) error
}
