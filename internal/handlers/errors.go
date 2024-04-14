package handlers

import "errors"

var ErrNoTag = errors.New("no tag provided")
var ErrIncorrectTag = errors.New("incorrect tag")
var ErrNoFeature = errors.New("no feature provided")
var ErrIncorrectFeature = errors.New("incorrect feature")
var ErrIncorrectLimit = errors.New("incorrect limit")
var ErrIncorrectOffset = errors.New("incorrect offset")
var ErrIncorrectID = errors.New("incorrect ID")
var ErrInternal = errors.New("internal server error")
var ErrIncorrectBannerContent = errors.New("incorrect banner content")
var ErrNoUsername = errors.New("no username provided")
var ErrNoPassword = errors.New("no password provided")
