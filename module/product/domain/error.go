package domain

import (
	"errors"
)

var (
	ErrProductNameCannotBeBlank = errors.New("product name cannot be blank")
)