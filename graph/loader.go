package graph

import (
	pathlib "pkb-agent/util/pathlib"
)

type Loader interface {
	Load(path pathlib.Path, callback func(node *Node) error) error
}

type GraphLoader struct {
	root       pathlib.Path
	nodeLoader Loader
}

func LoadGraph(root pathlib.Path, loader Loader) (*Graph, error) {
	graphLoader := GraphLoader{
		root:       root,
		nodeLoader: loader,
	}
	return graphLoader.Load()
}

func (gl *GraphLoader) Load() (*Graph, error) {
	builder, err := gl.LoadNodes()
	if err != nil {
		return nil, err
	}

	return builder.Finish()
}

func (gl *GraphLoader) LoadNodes() (*Builder, error) {
	builder := NewBuilder()

	callback := func(node *Node) error {
		return builder.AddNode(node)
	}

	if err := gl.nodeLoader.Load(gl.root, callback); err != nil {
		return nil, err
	}

	return builder, nil
}
