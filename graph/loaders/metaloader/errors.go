package metaloader

import "fmt"

type ErrUnknownLoader struct {
	unknownLoaderName string
}

func (err *ErrUnknownLoader) Error() string {
	return fmt.Sprintf("unknown loader %s", err.unknownLoaderName)
}

type ErrDuplicateLoaderName struct {
	duplicateLoaderName string
}

func (err *ErrDuplicateLoaderName) Error() string {
	return fmt.Sprintf("duplicate loader name %s", err.duplicateLoaderName)
}
