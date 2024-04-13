package usecase

import (
	"banners/internal/models"
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
