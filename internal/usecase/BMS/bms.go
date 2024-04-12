package BMS

import (
	"banners/internal/models"
	"context"
)

// Dependencies interfaces

type (
	Authentificator interface {
		auth(ctx context.Context, token string) (string, error)
	}
	BMSRepository interface {
		GetBanner(ctx context.Context, featureID models.FeatureID, tagID models.TagID) (models.Banner, error)
		SelectBanners(ctx context.Context, limit int64, offset int64) ([]models.Banner, error)
		AddBanner(ctx context.Context, banner models.Banner) error
		PatchBanner(ctx context.Context, banner models.Banner) error
		DeleteBanner(ctx context.Context, featureID models.FeatureID, tagID models.TagID) (models.Banner, error)
	}
)
