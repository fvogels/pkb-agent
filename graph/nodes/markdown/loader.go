package markdown

import (
	"fmt"
	"log/slog"
	"pkb-agent/graph"
	pathlib "pkb-agent/util/pathlib"
	"pkb-agent/util/sectionedfile"
	"strings"

	"gopkg.in/yaml.v3"
)

type Loader struct{}

func NewLoader() *Loader {
	return &Loader{}
}

type metadata struct {
	Name  string   `yaml:"name"`  // Name of the snippet.
	Links []string `yaml:"links"` // Links to other nodes.
}

func (loader *Loader) Load(path pathlib.Path, callback func(node *graph.Node) error) error {
	slog.Debug(
		"Loading node file",
		slog.String("loader", "markdown"),
		slog.String("path", path.String()),
	)

	unparsedMetadata, err := loader.extractMetadata(path)
	if err != nil {
		return fmt.Errorf("failed to extract metadata from %s: %v", path.String(), err)
	}

	metadata, err := loader.parseMetadata(unparsedMetadata)
	if err != nil {
		return fmt.Errorf("failed to parse metadata from %s: %v", path.String(), err)
	}

	node := graph.Node{
		Name:      metadata.Name,
		Type:      "markdown",
		Links:     append(metadata.Links, "Markdown"),
		Backlinks: nil,
		Extra: &Extra{
			Path: path,
		},
	}

	if err := callback(&node); err != nil {
		return err
	}

	return nil
}

func (loader *Loader) extractMetadata(path pathlib.Path) (string, error) {
	file, err := sectionedfile.LoadSectionedFile(path, isDelimiter)
	if err != nil {
		return "", err
	}

	if len(file.Sections) < 2 {
		return "", &ErrMalformed{path: path}
	}

	return strings.Join(file.Sections[1].Lines, "\n"), nil
}

func (loader *Loader) parseMetadata(unparsedMetadata string) (metadata, error) {
	var result metadata

	if err := yaml.Unmarshal([]byte(unparsedMetadata), &result); err != nil {
		return metadata{}, err
	}

	if len(result.Name) == 0 {
		slog.Error("Markdown node is missing name", slog.String("metadata", unparsedMetadata))
		return metadata{}, &ErrMissingName{}
	}

	return result, nil
}

func isDelimiter(line string) bool {
	return strings.TrimSpace(line) == "---"
}
