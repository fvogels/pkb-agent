package backblaze

import (
	"errors"
)

var ErrBucketMissing = errors.New("missing bucket")
var ErrFilenameMissing = errors.New("missing filename")
var ErrNameMissing = errors.New("missing node name")
