package model

import (
	"emperror.dev/errors"
)

const EntityNotFound = "entity not found"

var ErrEntityNotFound = errors.New("entity not found")
