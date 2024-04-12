package models

type Banner struct {
	FeatureID int32
	TagIDs    []int32
	Contents  BannerContents
	IsActive  bool
}
