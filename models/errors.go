package models

import "errors"

var ErrBothEmailPhoneRequest = errors.New("request contains both phone and email field is forbidden, use either phone or email")
var ErrEmptyParams = errors.New("request doesn't contains neither phone nor email, you must provide one")
