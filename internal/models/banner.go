package models

import "time"

type Banner struct {
	FeatureID FeatureID
	TagID     TagID
	PatchDate time.Time
	Contents  BannerContents
}
