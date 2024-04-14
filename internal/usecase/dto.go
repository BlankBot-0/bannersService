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

type BannerVersionsParams struct {
	TagID     int32  `json:"tag_id"`
	FeatureID int32  `json:"feature_id"`
	Limit     *int32 `json:"limit,omitempty"`
	Offset    *int32 `json:"offset,omitempty"`
}

type BannerDTO struct {
	ID        int64     `json:"banner_id"`
	FeatureID int32     `json:"feature_id"`
	Contents  string    `json:"content"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BannerVersionDTO struct {
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
}
type BannerVersionsDTO struct {
	ID       int32              `json:"banner_id"`
	Versions []BannerVersionDTO `json:"versions"`
}

func NewBannerVersionsDTO(rows []banners.ListBannerVersionsRow, bannerID int32) BannerVersionsDTO {
	versions := make([]BannerVersionDTO, len(rows))
	for i, row := range rows {
		versions[i] = BannerVersionDTO{
			Content:   string(row.Contents),
			UpdatedAt: row.UpdatedAt.Time,
		}
	}
	return BannerVersionsDTO{
		ID:       bannerID,
		Versions: versions,
	}
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
