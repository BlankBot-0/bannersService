package BMS

import (
	"banners/internal/models"
	"banners/internal/repository/postgres/banners"
	"banners/internal/usecase"
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
)

// Dependencies interfaces

type (
	BMSRepository interface {
		AddBannerTags(ctx context.Context, arg banners.AddBannerTagsParams) error
		CheckActiveUserBanner(ctx context.Context, arg banners.CheckActiveUserBannerParams) (bool, error)
		CheckExistsBanner(ctx context.Context, arg banners.CheckExistsBannerParams) (int64, error)
		CheckBannerId(ctx context.Context, bannerID int32) (bool, error)
		CreateBanner(ctx context.Context, arg banners.CreateBannerParams) (int64, error)
		CreateBannerInfo(ctx context.Context, arg banners.CreateBannerInfoParams) error
		DeleteBannerInfo(ctx context.Context, bannerID int32) error
		DeleteBannerTags(ctx context.Context, bannerID int32) error
		DeleteBanner(ctx context.Context, bannerID int32) error
		GetUserBanner(ctx context.Context, arg banners.GetUserBannerParams) (banners.GetUserBannerRow, error)
		ListBannerVersions(ctx context.Context, arg banners.ListBannerVersionsParams) ([]banners.ListBannerVersionsRow, error)
		ListBanners(ctx context.Context, arg banners.ListBannersParams) ([]banners.ListBannersRow, error)
		UpdateBannerContents(ctx context.Context, arg banners.UpdateBannerContentsParams) error
		UpdateBannerFeature(ctx context.Context, arg banners.UpdateBannerFeatureParams) error
		UpdateBannerIsActive(ctx context.Context, arg banners.UpdateBannerIsActiveParams) error
		WithTx(tx pgx.Tx) *banners.Queries
	}
	txBuilder interface {
		Begin(ctx context.Context) (pgx.Tx, error)
	}
	db interface {
		banners.DBTX
		txBuilder
	}
)

type Deps struct {
	Repository BMSRepository
	TxBuilder  db
}

type BMS struct {
	Deps
}

func NewBMS(deps Deps) *BMS {
	return &BMS{
		Deps: deps,
	}
}

func (s *BMS) CreateBanner(ctx context.Context, banner usecase.BannerJsonDTO) error {
	tx, err := s.TxBuilder.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := s.Repository.WithTx(tx)

	_, err = qtx.CheckExistsBanner(ctx, banners.CheckExistsBannerParams{
		FeatureID: banner.FeatureID,
		TagIds:    banner.BannerWithTagsDTO.Tags,
	})
	if !errors.Is(err, pgx.ErrNoRows) {
		return ErrFeatureTagPairAlreadyExists
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	bannerId, err := qtx.CreateBanner(ctx, banners.CreateBannerParams{
		FeatureID: banner.FeatureID,
		IsActive:  banner.IsActive,
	})
	if err != nil {
		return err
	}

	bannerContent, err := json.Marshal(banner.Contents)
	if err != nil {
		return err
	}
	err = qtx.CreateBannerInfo(ctx, banners.CreateBannerInfoParams{
		Contents: bannerContent,
		BannerID: int32(bannerId),
	})
	if err != nil {
		return err
	}

	err = qtx.AddBannerTags(ctx, banners.AddBannerTagsParams{
		BannerID: int32(bannerId),
		TagIds:   banner.BannerWithTagsDTO.Tags,
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *BMS) UserBanner(ctx context.Context, tagID int32, featureID int32) (models.BannerContent, error) {
	tx, err := s.TxBuilder.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)
	qtx := s.Repository.WithTx(tx)

	banner, err := qtx.GetUserBanner(ctx, banners.GetUserBannerParams{
		TagID:     tagID,
		FeatureID: featureID,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrBannerNotFound
		}
		return "", err
	}
	if !banner.IsActive {
		return "", ErrNotActiveBanner
	}
	err = tx.Commit(ctx)
	if err != nil {
		return "", err
	}

	return models.BannerContent(banner.Contents), nil
}

func (s *BMS) ListBanners(ctx context.Context, arg banners.ListBannersParams) ([]banners.ListBannersRow, error) {
	tx, err := s.TxBuilder.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	qtx := s.Repository.WithTx(tx)

	Banners, err := qtx.ListBanners(ctx, banners.ListBannersParams{
		TagID:     arg.TagID,
		FeatureID: arg.FeatureID,
		OffsetVal: arg.OffsetVal,
		LimitVal:  arg.LimitVal,
	})
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return Banners, nil
}

func (s *BMS) ListBannerVersions(ctx context.Context, arg usecase.BannerVersionsParams) (usecase.BannerVersionsDTO, error) {
	tx, err := s.TxBuilder.Begin(ctx)
	if err != nil {
		return usecase.BannerVersionsDTO{}, err
	}
	defer tx.Rollback(ctx)
	qtx := s.Repository.WithTx(tx)

	bannerID, err := qtx.CheckExistsBanner(ctx, banners.CheckExistsBannerParams{
		FeatureID: arg.FeatureID,
		TagIds:    []int32{arg.TagID},
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return usecase.BannerVersionsDTO{}, ErrBannerNotFound
	}
	if err != nil {
		return usecase.BannerVersionsDTO{}, err
	}

	limit := DefaultBannerVersionsLimit
	if arg.Limit != nil {
		limit = *arg.Limit
	}
	offset := DefaultBannerVersionsOffset
	if arg.Offset != nil {
		offset = *arg.Offset
	}
	versions, err := qtx.ListBannerVersions(ctx, banners.ListBannerVersionsParams{
		BannerID:  int32(bannerID),
		OffsetVal: offset,
		LimitVal:  limit,
	})
	if err != nil {
		return usecase.BannerVersionsDTO{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return usecase.BannerVersionsDTO{}, err
	}

	return usecase.NewBannerVersionsDTO(versions, int32(bannerID)), nil

}

func (s *BMS) UpdateBanner(ctx context.Context, params usecase.UpdateBannerDTO) error {
	tx, err := s.TxBuilder.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := s.Repository.WithTx(tx)

	bannerID := params.BannerID
	existsBanner, err := qtx.CheckBannerId(ctx, bannerID)
	if err != nil {
		return err
	}
	if !existsBanner {
		return ErrBannerIDNotFound
	}

	if params.FeatureID != nil {
		featureID := *params.FeatureID
		err = qtx.UpdateBannerFeature(ctx, banners.UpdateBannerFeatureParams{
			FeatureID: featureID,
			BannerID:  bannerID,
		})
		if err != nil {
			return err
		}
	}

	if params.IsActive != nil {
		isActive := *params.IsActive
		err = qtx.UpdateBannerIsActive(ctx, banners.UpdateBannerIsActiveParams{
			IsActive: isActive,
			BannerID: bannerID,
		})
		if err != nil {
			return err
		}
	}

	if params.TagIDs != nil {
		tags := params.TagIDs
		err = qtx.DeleteBannerTags(ctx, bannerID)
		if err != nil {
			return err
		}

		err = qtx.AddBannerTags(ctx, banners.AddBannerTagsParams{
			BannerID: bannerID,
			TagIds:   tags,
		})
		if err != nil {
			return err
		}
	}

	if params.Content != nil {
		bannerContent, err := json.Marshal(params.Content)
		if err != nil {
			return err
		}
		err = qtx.UpdateBannerContents(ctx, banners.UpdateBannerContentsParams{
			BannerID: bannerID,
			Contents: bannerContent,
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *BMS) DeleteBanner(ctx context.Context, id int32) error {
	tx, err := s.TxBuilder.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := s.Repository.WithTx(tx)

	existsBanner, err := qtx.CheckBannerId(ctx, id)
	if err != nil {
		return err
	}
	if !existsBanner {
		return ErrBannerNotFound
	}

	err = qtx.DeleteBannerTags(ctx, id)
	if err != nil {
		return err
	}

	err = qtx.DeleteBannerInfo(ctx, id)
	if err != nil {
		return err
	}

	err = qtx.DeleteBanner(ctx, id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

const DefaultBannerVersionsLimit = int32(3)
const DefaultBannerVersionsOffset = int32(0)
