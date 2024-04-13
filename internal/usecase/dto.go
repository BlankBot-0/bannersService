package usecase

import (
	"banners/internal/models"
	"banners/internal/repository/postgres/banners"
	"time"
)

type UpdateBannerParams struct {
	BannerID  int32                 `json:"banner_id"`
	TagIDs    []int32               `json:"tag_ids,omitempty"`
	FeatureID int32                 `json:"feature_id,omitempty"`
	IsActive  bool                  `json:"is_active,omitempty"`
	Contents  models.BannerContents `json:"contents,omitempty"`
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

func (p UpdateBannerParams) BannerContents() (models.BannerContents, bool) {
	_, ok := p.Paramsmap["banner_contents"]
	if !ok {
		return nil, false
	}
	return p.Contents, true
}

type UserBannerParams struct {
	TagID     int32 `json:"tag_id"`
	FeatureID int32 `json:"feature_id"`
}

type ListBannerVersionsResponse struct {
	Versions []ListBannerVersionsResponseRow `json:"version"`
}

type ListBannerVersionsResponseRow struct {
	ID        int64     `json:"id"`
	FeatureID int32     `json:"feature_id"`
	Contents  []byte    `json:"contents"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewListBannerVersionsResponse(rows []banners.ListBannerVersionsRow) ListBannerVersionsResponse {
	versions := make([]ListBannerVersionsResponseRow, len(rows))
	for i, row := range rows {
		versions[i] = ListBannerVersionsResponseRow{
			ID:        row.ID,
			FeatureID: row.FeatureID,
			Contents:  row.Contents,
			IsActive:  row.IsActive,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		}
	}
	return ListBannerVersionsResponse{
		Versions: versions,
	}
}

type ListBannersResponse struct {
	Banners []ListBannersResponseRow `json:"banners"`
}

type ListBannersResponseRow struct {
	ID        int64     `json:"id"`
	FeatureID int32     `json:"feature_id"`
	Contents  []byte    `json:"contents"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Tags      []int32   `json:"tags"`
}

func NewListBannersResponse(rows []banners.ListBannersRow) ListBannersResponse {
	bannersList := make([]ListBannersResponseRow, len(rows))
	for i, row := range rows {
		bannersList[i] = ListBannersResponseRow{
			ID:        row.ID,
			FeatureID: row.FeatureID,
			Contents:  row.Contents,
			IsActive:  row.IsActive,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
			Tags:      row.Tags,
		}
	}
	return ListBannersResponse{
		Banners: bannersList,
	}
}
