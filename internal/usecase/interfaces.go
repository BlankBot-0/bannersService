package usecase

import (
	"banners/internal/models"
	"banners/internal/repository/postgres/banners"
	"context"
)

type BannerManagementSystem interface {
	CreateBanner(ctx context.Context, banner BannerDTO) error
	UserBanner(ctx context.Context, tagID, featureID int32) (models.BannerContent, error)
	ListBanners(ctx context.Context, arg banners.ListBannersParams) ([]banners.ListBannersRow, error)
	ListBannerVersions(ctx context.Context, arg BannerVersionsParams) (BannerVersionsDTO, error)
	UpdateBanner(ctx context.Context, params UpdateBannerDTO) error
	DeleteBanner(ctx context.Context, id int32) error
}

type AuthenticationSystem interface {
	UserToken(ctx context.Context, credentials CredentialsDTO) (string, error)
	UserAuth(ctx context.Context, token string) error
	AdminToken(ctx context.Context, credentials CredentialsDTO) (string, error)
	AdminAuth(ctx context.Context, token string) error
}

type Cacher interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
}
