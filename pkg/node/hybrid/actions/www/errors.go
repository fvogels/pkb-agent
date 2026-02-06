package www

import "errors"

var (
	ErrMissingAttribute = errors.New("download action is missing mandatory attribute")
)
