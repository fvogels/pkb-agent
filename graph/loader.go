package graph

import (
	"fmt"
	pathlib "pkb-agent/util/pathlib"
)

type Loader interface {
	Load(path pathlib.Path, callback func(node *Node) error) error
}

func LoadGraph(root pathlib.Path, loader Loader) error {
	callback := func(entry *Node) error {
		fmt.Printf("%s\n", entry.Name)
		return nil
	}

	if err := loader.Load(root, callback); err != nil {
		return err
	}

	return nil
}
