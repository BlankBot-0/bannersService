package models

type Banner struct {
	TagIDs    []int32       `json:"tag_ids"`
	FeatureID int32         `json:"feature_id"`
	Content   BannerContent `json:"content"`
	IsActive  bool          `json:"is_active"`
}
