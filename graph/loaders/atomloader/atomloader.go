package atomloader

import (
	"errors"
	"pkb-agent/graph"
	"pkb-agent/graph/nodes/atom"
	pathlib "pkb-agent/util/path"

	"gopkg.in/yaml.v2"
)

type Loader struct{}

func New() *Loader {
	return &Loader{}
}

type entry struct {
	Name       string
	Identifier string
	Links      []string
}

func (loader *Loader) Load(path pathlib.Path, callback func(node graph.Node) error) error {
	source, err := path.ReadFile()
	if err != nil {
		return err
	}

	var entries []entry
	if err := yaml.Unmarshal(source, &entries); err != nil {
		return err
	}

	var errs []error
	for _, entry := range entries {
		node := atom.Node{
			Name:       entry.Name,
			Identifier: entry.Identifier,
			Links:      entry.Links,
		}

		if err := callback(&node); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
