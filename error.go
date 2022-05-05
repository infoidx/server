package server

import "errors"

var (
	ErrGinInstanceNotInit = errors.New("the gin instance has not been initialized")
)
