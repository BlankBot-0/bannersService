package BMS

import (
	"banners/internal/models"
	"banners/internal/repository/postgres/banners"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Dependencies interfaces

type (
	BMSRepository interface {
		AddBannerTags(ctx context.Context, arg banners.AddBannerTagsParams) error
		CheckActiveUserBanner(ctx context.Context, arg banners.CheckActiveUserBannerParams) (bool, error)
		CheckExistsBanner(ctx context.Context, arg banners.CheckExistsBannerParams) (bool, error)
		CheckBannerId(ctx context.Context, id int64) (bool, error)
		CreateBanner(ctx context.Context, arg banners.CreateBannerParams) (int64, error)
		CreateBannerInfo(ctx context.Context, arg banners.CreateBannerInfoParams) error
		DeleteBannerInfo(ctx context.Context, bannerID pgtype.Int4) error
		DeleteBannerTags(ctx context.Context, bannerID pgtype.Int4) error
		GetUserBanner(ctx context.Context, arg banners.GetUserBannerParams) ([]byte, error)
		ListBannerVersions(ctx context.Context, arg banners.ListBannerVersionsParams) ([]banners.ListBannerVersionsRow, error)
		ListBanners(ctx context.Context, arg banners.ListBannersParams) ([]banners.ListBannersRow, error)
		UpdateBannerContents(ctx context.Context, arg banners.UpdateBannerContentsParams) error
		UpdateBannerFeature(ctx context.Context, arg banners.UpdateBannerFeatureParams) error
		UpdateBannerIsActive(ctx context.Context, arg banners.UpdateBannerIsActiveParams) error
		WithTx(tx pgx.Tx) *banners.Queries
	}
)

type Deps struct {
	Repository BMSRepository
}

type BMS struct {
	Deps
}

func NewBMS(deps Deps) *BMS {
	return &BMS{
		Deps: deps,
	}
}

func (s *BMS) CreateBanner(ctx context.Context, banner models.Banner) error {
	existsBanner, err := s.Repository.CheckExistsBanner(ctx, banners.CheckExistsBannerParams{
		FeatureID: banner.FeatureID,
		TagIds:    banner.TagIDs,
	})
	if err != nil {
		return err
	}
	if existsBanner {
		return ErrFeatureTagPairAlreadyExists
	}

	bannerId, err := s.Repository.CreateBanner(ctx, banners.CreateBannerParams{
		FeatureID: banner.FeatureID,
		IsActive:  banner.IsActive,
	})
	if err != nil {
		return err
	}

	err = s.Repository.CreateBannerInfo(ctx, banners.CreateBannerInfoParams{
		Contents: banner.Contents,
		BannerID: int32(bannerId),
	})
	if err != nil {
		return err
	}

	err = s.Repository.AddBannerTags(ctx, banners.AddBannerTagsParams{
		BannerID: int32(bannerId),
		TagIds:   banner.TagIDs,
	})
	if err != nil {
		return err
	}
	return nil
}
