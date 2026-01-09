package pkg

import (
	"errors"
)

var ErrNameClash = errors.New("nodes share name")
var ErrUnknownNode = errors.New("unknown node")
