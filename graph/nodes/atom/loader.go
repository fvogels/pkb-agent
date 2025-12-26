package atom

import (
	"errors"
	"log/slog"
	"pkb-agent/graph"
	"pkb-agent/util/pathlib"

	"gopkg.in/yaml.v3"
)

type Loader struct{}

func NewLoader() *Loader {
	return &Loader{}
}

type entry struct {
	Name  string
	Links []string
}

func New(atomName string) *graph.Node {
	return &graph.Node{
		Name:      atomName,
		Type:      "atom",
		Links:     nil,
		Backlinks: nil,
		Info:      &Info{},
	}
}

func (loader *Loader) Load(path pathlib.Path, callback func(node *graph.Node) error) error {
	slog.Debug(
		"Loading node file",
		slog.String("loader", "atom"),
		slog.String("path", path.String()),
	)

	info := Info{}

	source, err := path.ReadFile()
	if err != nil {
		return err
	}

	var entries []entry
	if err := yaml.Unmarshal(source, &entries); err != nil {
		return err
	}

	var errs []error
	for index, entry := range entries {
		if len(entry.Name) == 0 {
			errs = append(errs, &ErrMissingName{path: path, index: index})
			continue
		}

		node := graph.Node{
			Name:      entry.Name,
			Type:      "atom",
			Links:     entry.Links,
			Backlinks: nil,
			Info:      &info,
		}

		if err := callback(&node); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
