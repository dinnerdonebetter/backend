package database

import (
	"errors"
)

var (
	ErrAlreadyFinalized = errors.New("meal plan already finalized")
)
