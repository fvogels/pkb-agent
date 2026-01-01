package graph

import (
	"pkb-agent/graph/node"
	pathlib "pkb-agent/util/pathlib"
)

type GraphLoader struct {
	root       pathlib.Path
	nodeLoader node.Loader
}

func LoadGraph(root pathlib.Path, loader node.Loader) (*Graph, error) {
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

	return builder.Finalize()
}

func (gl *GraphLoader) LoadNodes() (*Builder, error) {
	builder := NewBuilder()

	callback := func(node node.RawNode) error {
		return builder.AddNode(node)
	}

	if err := gl.nodeLoader.Load(gl.root, callback); err != nil {
		return nil, err
	}

	return builder, nil
}
