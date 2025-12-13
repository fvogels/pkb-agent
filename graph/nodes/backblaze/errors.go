package backblaze

import (
	"fmt"
	"pkb-agent/util/pathlib"
)

type ErrBucketMissingName struct {
	path  pathlib.Path
	index int
}

func (err *ErrBucketMissingName) Error() string {
	return fmt.Sprintf("bucket missing name; file %s, index %d", err.path.String(), err.index)
}

type ErrFileMissingName struct {
	bucket string
	path   pathlib.Path
	index  int
}

func (err *ErrFileMissingName) Error() string {
	return fmt.Sprintf("backblaze node missing name in file %s, bucket %s, index %d", err.path.String(), err.bucket, err.index)
}
