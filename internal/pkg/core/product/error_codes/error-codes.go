package error_codes

import "github.com/pkg/errors"

var (
	ErrExternalProblem   = errors.New("Some problem occured")
	ErrNameTooShortError = errors.New("name is too short")
	ErrNegativePrice     = errors.New("price should be positive")
	ErrProductNotExist   = errors.New("product does not exist")
	ErrProductExists     = errors.New("product exists")
)

func GetInternal(err error) error {
	if errors.Is(err, ErrProductNotExist) {
		return ErrProductNotExist
	}
	if errors.Is(err, ErrProductExists) {
		return ErrProductExists
	}
	if errors.Is(err, ErrNegativePrice) {
		return ErrNegativePrice
	}
	if errors.Is(err, ErrNameTooShortError) {
		return ErrNameTooShortError
	}
	return ErrExternalProblem
}
