package usecase

import (
	"banners/internal/models"
	"banners/internal/repository/postgres/banners"
	"time"
)

type UpdateBannerParams struct {
	BannerID  int32                `json:"banner_id"`
	TagIDs    []int32              `json:"tag_ids,omitempty"`
	FeatureID int32                `json:"feature_id,omitempty"`
	IsActive  bool                 `json:"is_active,omitempty"`
	Content   models.BannerContent `json:"contents,omitempty"`
	Paramsmap map[string]bool
}

func (p UpdateBannerParams) Tags() ([]int32, bool) {
	_, ok := p.Paramsmap["tag_ids"]
	if !ok {
		return nil, false
	}
	return p.TagIDs, true
}

func (p UpdateBannerParams) Feature() (int32, bool) {
	_, ok := p.Paramsmap["feature_id"]
	if !ok {
		return 0, false
	}
	return p.FeatureID, true
}

func (p UpdateBannerParams) Active() (bool, bool) {
	_, ok := p.Paramsmap["is_active"]
	if !ok {
		return false, false
	}
	return p.IsActive, true
}

func (p UpdateBannerParams) BannerContent() (models.BannerContent, bool) {
	_, ok := p.Paramsmap["banner_contents"]
	if !ok {
		return "", false
	}
	return p.Content, true
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
