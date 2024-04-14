package usecase

import (
	"banners/internal/models"
	"banners/internal/repository/postgres/banners"
	"context"
)

type BannerManagementSystem interface {
	CreateBanner(ctx context.Context, banner BannerJsonDTO) error
	UserBanner(ctx context.Context, tagID, featureID int32) (models.BannerContent, error)
	ListBanners(ctx context.Context, arg banners.ListBannersParams) ([]banners.ListBannersRow, error)
	UpdateBanner(ctx context.Context, params UpdateBannerDTO) error
	DeleteBanner(ctx context.Context, id int32) error
}
