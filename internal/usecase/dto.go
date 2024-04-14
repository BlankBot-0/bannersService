package usecase

import (
	"banners/internal/repository/postgres/banners"
	"time"
)

type UpdateBannerDTO struct {
	BannerID  int32
	TagIDs    []int32           `json:"tag_ids,omitempty"`
	FeatureID *int32            `json:"feature_id,omitempty"`
	IsActive  *bool             `json:"is_active,omitempty"`
	Content   map[string]string `json:"content,omitempty"`
}

type UserBannerParams struct {
	TagID     int32 `json:"tag_id"`
	FeatureID int32 `json:"feature_id"`
}

type BannerDTO struct {
	ID        int64     `json:"banner_id"`
	FeatureID int32     `json:"feature_id"`
	Contents  string    `json:"content"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewBannerVersionsDTO(rows []banners.ListBannerVersionsRow) []BannerDTO {
	versions := make([]BannerDTO, len(rows))
	for i, row := range rows {
		versions[i] = BannerDTO{
			ID:        row.ID,
			FeatureID: row.FeatureID,
			Contents:  string(row.Contents),
			IsActive:  row.IsActive,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		}
	}
	return versions
}

type BannerWithTagsDTO struct {
	BannerDTO
	Tags []int32 `json:"tag_ids"`
}

func NewListBannersDTO(rows []banners.ListBannersRow) []BannerWithTagsDTO {
	bannersList := make([]BannerWithTagsDTO, len(rows))
	for i, row := range rows {
		bannersList[i] = BannerWithTagsDTO{
			BannerDTO: BannerDTO{
				ID:        row.ID,
				FeatureID: row.FeatureID,
				Contents:  string(row.Contents),
				IsActive:  row.IsActive,
				CreatedAt: row.CreatedAt.Time,
				UpdatedAt: row.UpdatedAt.Time,
			},
			Tags: row.Tags,
		}
	}
	return bannersList
}

type BannerJsonDTO struct {
	BannerWithTagsDTO
	Contents map[string]string `json:"content"`
}
