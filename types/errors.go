package types

import "errors"

var (
	ErrInvalidEntityType   = errors.New("depot: invalid entity type")
	ErrInvalidTransform    = errors.New("depot: invalid transform")
	ErrEntityNotFound      = errors.New("depot: entity not found")
	ErrEntityAlreadyExists = errors.New("depot: entity already exists")
)
