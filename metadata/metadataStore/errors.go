package metadatastore

import "errors"

var (
	ErrInvalidName    = errors.New("invalid company  or username")
	ErrInvalidPasswd  = errors.New("invalid password")
	ErrDbErrorGeneric = errors.New("database error: ")
	ErrCompanyNoExist = errors.New("comapny doesnt exist")
	ErrDeviceNotExist = errors.New("device doesnt exist")
)
