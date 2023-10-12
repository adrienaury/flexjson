package flexjson

import "errors"

var (
	ErrUnexpectedToken = errors.New("unexpected token")
	ErrInvalidKeyType  = errors.New("invalid key type")
)
