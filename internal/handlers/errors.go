package handlers

import "errors"

var ErrNoTag = errors.New("no tag provided")
var ErrIncorrectTag = errors.New("incorrect tag")
var ErrNoFeature = errors.New("no feature provided")
var ErrIncorrectFeature = errors.New("incorrect feature")
