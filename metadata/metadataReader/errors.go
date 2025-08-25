package metadatareader

import "errors"

var (
	ErrComanyNotFound = errors.New("company not found")
	ErrDbErrorGeneric = errors.New("database error: ")
)
