package bookmark

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
	Name        string   `yaml:"name"`
	Links       []string `yaml:"links"`
	URL         string   `yaml:"url"`
	Description string   `yaml:"description"`
}

func (loader *Loader) Load(path pathlib.Path, callback func(node *graph.Node) error) error {
	slog.Debug(
		"Loading node file",
		slog.String("loader", "bookmark"),
		slog.String("path", path.String()),
	)

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
			slog.Debug(
				"Missing name",
				slog.String("path", path.String()),
				slog.Int("index", index),
			)
			errs = append(errs, &ErrMissingName{path: path, index: index})
			continue
		}

		if len(entry.Description) == 0 {
			slog.Debug(
				"Missing description",
				slog.String("path", path.String()),
				slog.Int("index", index),
			)

			errs = append(errs, &ErrMissingDescription{path: path, index: index})
			continue
		}

		node := graph.Node{
			Name:      entry.Name,
			Type:      "bookmark",
			Links:     append(entry.Links, "Bookmark"),
			Backlinks: nil,
			Info: &Info{
				URL:         entry.URL,
				Description: entry.Description,
			},
		}

		if err := callback(&node); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
