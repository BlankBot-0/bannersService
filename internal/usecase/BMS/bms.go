package BMS

import (
	"banners/internal/models"
	"banners/internal/repository/postgres/banners"
	"banners/internal/usecase"
	"context"
	"github.com/jackc/pgx/v5"
)

// Dependencies interfaces

type (
	BMSRepository interface {
		AddBannerTags(ctx context.Context, arg banners.AddBannerTagsParams) error
		CheckActiveUserBanner(ctx context.Context, arg banners.CheckActiveUserBannerParams) (bool, error)
		CheckExistsBanner(ctx context.Context, arg banners.CheckExistsBannerParams) (bool, error)
		CheckBannerId(ctx context.Context, bannerID int32) (bool, error)
		CreateBanner(ctx context.Context, arg banners.CreateBannerParams) (int64, error)
		CreateBannerInfo(ctx context.Context, arg banners.CreateBannerInfoParams) error
		DeleteBannerInfo(ctx context.Context, bannerID int32) error
		DeleteBannerTags(ctx context.Context, bannerID int32) error
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

//TODO: wrap every query within method in transaction

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
func (s *BMS) UserBanner(ctx context.Context, tagID int32, featureID int32) (models.BannerContents, error) {
	existsBanner, err := s.Repository.CheckExistsBanner(ctx, banners.CheckExistsBannerParams{
		FeatureID: featureID,
		TagIds:    []int32{tagID},
	})
	if err != nil {
		return nil, err
	}
	if !existsBanner {
		return nil, ErrBannerNotFound
	}

	activeBanner, err := s.Repository.CheckActiveUserBanner(ctx, banners.CheckActiveUserBannerParams{
		TagID:     tagID,
		FeatureID: featureID,
	})
	if err != nil {
		return nil, err
	}
	if !activeBanner {
		return nil, ErrNotActiveBanner
	}

	banner, err := s.Repository.GetUserBanner(ctx, banners.GetUserBannerParams{
		TagID:     tagID,
		FeatureID: featureID,
	})
	if err != nil {
		return nil, err
	}
	return banner, nil
}

func (s *BMS) ListBanners(ctx context.Context, arg banners.ListBannersParams) ([]banners.ListBannersRow, error) {
	Banners, err := s.Repository.ListBanners(ctx, banners.ListBannersParams{
		TagID:     arg.TagID,
		FeatureID: arg.FeatureID,
		OffsetVal: arg.OffsetVal,
		LimitVal:  arg.LimitVal,
	})
	if err != nil {
		return nil, err
	}
	return Banners, nil
}

func (s *BMS) UpdateBanner(ctx context.Context, params usecase.UpdateBannerParams) error {
	bannerID := params.BannerID
	existsBanner, err := s.Repository.CheckBannerId(ctx, bannerID)
	if err != nil {
		return err
	}
	if !existsBanner {
		return ErrBannerNotFound
	}

	featureID, ok := params.Feature()
	if ok {
		err = s.Repository.UpdateBannerFeature(ctx, banners.UpdateBannerFeatureParams{
			FeatureID: featureID,
			BannerID:  bannerID,
		})
		if err != nil {
			return err
		}
	}

	isActive, ok := params.Active()
	if ok {
		err = s.Repository.UpdateBannerIsActive(ctx, banners.UpdateBannerIsActiveParams{
			IsActive: isActive,
			BannerID: bannerID,
		})
		if err != nil {
			return err
		}
	}

	tags, ok := params.Tags()
	if ok {
		err = s.Repository.DeleteBannerTags(ctx, bannerID)
		if err != nil {
			return err
		}

		err = s.Repository.AddBannerTags(ctx, banners.AddBannerTagsParams{
			BannerID: bannerID,
			TagIds:   tags,
		})
		if err != nil {
			return err
		}
	}

	contents, ok := params.BannerContents()
	if ok {
		err = s.Repository.UpdateBannerContents(ctx, banners.UpdateBannerContentsParams{
			BannerID: bannerID,
			Contents: contents,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *BMS) DeleteBanner(ctx context.Context, id int32) error {
	existsBanner, err := s.Repository.CheckBannerId(ctx, id)
	if err != nil {
		return err
	}
	if !existsBanner {
		return ErrBannerNotFound
	}

	err = s.Repository.DeleteBannerTags(ctx, id)
	if err != nil {
		return err
	}

	err = s.Repository.DeleteBannerInfo(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
