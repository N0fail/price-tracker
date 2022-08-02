package error_codes

import "github.com/pkg/errors"

var (
	ErrNameTooShortError = errors.New("name is too short")
	ErrNegativePrice     = errors.New("price should be positive")
	ErrProductNotExist   = errors.New("product does not exist")
	ErrProductExists     = errors.New("product exists")
)
