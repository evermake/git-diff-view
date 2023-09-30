package diff

import (
	"errors"
	"fmt"
)

var (
	ErrMalformedDiff       = errors.New("malformed diff")
	ErrUnknownStatusLetter = fmt.Errorf("%w: invalid status letter", ErrMalformedDiff)
)
