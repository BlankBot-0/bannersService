package handlers

import "errors"

var ErrNoTag = errors.New("no tag provided")
var ErrIncorrectTag = errors.New("incorrect tag")
var ErrNoFeature = errors.New("no feature provided")
var ErrIncorrectFeature = errors.New("incorrect feature")
var ErrIncorrectLimit = errors.New("incorrect limit")
var ErrIncorrectOffset = errors.New("incorrect offset")
var ErrInternal = errors.New("internal server error")
var ErrIncorrectBannerContent = errors.New("incorrect banner content")
