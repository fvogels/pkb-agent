package pkg

import (
	"pkb-agent/pkg/node"
	pathlib "pkb-agent/util/pathlib"

	"gopkg.in/yaml.v3"
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
	configuration, err := gl.loadRawConfiguration()
	if err != nil {
		return nil, err
	}

	builder := NewBuilder()

	callback := func(node node.RawNode) error {
		return builder.AddNode(node)
	}

	if err := gl.nodeLoader.Load(gl.root.Parent(), configuration, callback); err != nil {
		return nil, err
	}

	return builder, nil
}

func (gl *GraphLoader) loadRawConfiguration() (any, error) {
	buffer, err := gl.root.ReadFile()
	if err != nil {
		return nil, err
	}

	var result any
	if err := yaml.Unmarshal(buffer, &result); err != nil {
		return nil, err
	}

	return result, nil
}
