package BMS

import "errors"

var ErrFeatureTagPairAlreadyExists = errors.New("feature tag pair already exists")
var ErrBannerNotFound = errors.New("banner with given feature and tag not found")
var ErrNotActiveBanner = errors.New("banner is not active")
