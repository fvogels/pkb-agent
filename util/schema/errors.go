package schema

import "errors"

var ErrNotAMap = errors.New("not a map")
var ErrMissingKey = errors.New("missing key")
var ErrWrongType = errors.New("wrong type")
