package download

import "errors"

var (
	ErrMissingAttribute = errors.New("download action is missing mandatory attribute")
	ErrUnknownSource    = errors.New("download action has unknown source")
)
