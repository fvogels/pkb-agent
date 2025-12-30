package bookmark

import (
	"errors"
)

var ErrMissingName = errors.New("bookmark node is missing name")
var ErrMissingDescription = errors.New("bookmark node is missing description")
var ErrMissingURL = errors.New("bookmark is missing an url")
