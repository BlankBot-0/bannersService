package usecase

import (
	"banners/internal/models"
	"context"
)

type BannerManagementSystem interface {
	Banner(ctx context.Context, featureID models.FeatureID, tagID models.TagID) (models.Banner, error)
	Banners(ctx context.Context, limit int64, offset int64) ([]models.Banner, error)
	NewBanner(ctx context.Context, banner models.Banner) error
	PatchBanner(ctx context.Context, banner models.Banner) error
	DeleteBanner(ctx context.Context, featureID models.FeatureID, tagID models.TagID) (models.Banner, error)
}
