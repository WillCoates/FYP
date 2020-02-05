package service

import "errors"

var ErrNoToken = errors.New("No token")
var ErrNoPermission = errors.New("Token lacking permission")
