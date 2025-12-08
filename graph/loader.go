package graph

import (
	"log/slog"
	pathlib "pkb-agent/util/pathlib"
)

type Loader interface {
	Load(path pathlib.Path, callback func(node *Node) error) error
}

type GraphLoader struct {
	root       pathlib.Path
	nodeLoader Loader
}

func LoadGraph(root pathlib.Path, loader Loader) error {
	graphLoader := GraphLoader{
		root:       root,
		nodeLoader: loader,
	}
	return graphLoader.Load()
}

func (gl *GraphLoader) Load() error {
	nodes, err := gl.LoadNodes()
	if err != nil {
		return err
	}

	if err := gl.EnsureLinkedNodeExistence(nodes); err != nil {
		return err
	}

	gl.AddBackLinks(nodes)

	return nil
}

type ErrNameClash struct{}

func (err *ErrNameClash) Error() string {
	return "name clash"
}

type ErrUnknownNodes struct{}

func (err *ErrUnknownNodes) Error() string {
	return "unknown node"
}

func (gl *GraphLoader) LoadNodes() (map[string]*Node, error) {
	nodes := make(map[string]*Node)
	duplicatesFound := false

	callback := func(node *Node) error {
		if _, alreadyExists := nodes[node.Name]; alreadyExists {
			slog.Debug("Multiple nodes with same name", slog.String("name", node.Name))
			duplicatesFound = true
		}

		nodes[node.Name] = node
		return nil
	}

	if err := gl.nodeLoader.Load(gl.root, callback); err != nil {
		return nil, err
	}

	if duplicatesFound {
		return nodes, &ErrNameClash{}
	}

	return nodes, nil
}

func (gl *GraphLoader) EnsureLinkedNodeExistence(nodes map[string]*Node) error {
	foundUnknownLinks := false

	for _, node := range nodes {
		for _, link := range node.Links {
			if _, found := nodes[link]; !found {
				slog.Debug("Unknown link", slog.String("node", node.Name), slog.String("link target", link))
				foundUnknownLinks = true
			}
		}
	}

	if foundUnknownLinks {
		return &ErrUnknownNodes{}
	}

	return nil
}

func (gl *GraphLoader) AddBackLinks(nodes map[string]*Node) {
	for _, node := range nodes {
		for _, link := range node.Links {
			linkedNode, ok := nodes[link]
			if !ok {
				panic("missing node; should have been noticed earlier")
			}

			linkedNode.Backlinks = append(linkedNode.Backlinks, node.Name)
		}
	}
}
