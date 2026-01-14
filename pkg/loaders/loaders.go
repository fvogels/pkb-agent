package loaders

import (
	"errors"
	"fmt"
	"pkb-agent/pkg/node"
)

var loaderTable map[string]node.Loader = make(map[string]node.Loader)

var ErrUnknownLoader = errors.New("unknown loader")

func RegisterLoader(id string, loader node.Loader) {
	loaderTable[id] = loader
}

func GetLoader(id string) (node.Loader, error) {
	loader, ok := loaderTable[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownLoader, id)
	}

	return loader, nil
}
