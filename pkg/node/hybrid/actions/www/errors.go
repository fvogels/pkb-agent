package www

import "errors"

var (
	ErrMissingDescription = errors.New("www action is missing description")
	ErrMissingURL         = errors.New("www action is missing url")
)
