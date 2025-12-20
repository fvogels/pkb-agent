package hybrid

import (
	"fmt"
	"log/slog"
	"pkb-agent/graph"
	"pkb-agent/util/multifile"
	pathlib "pkb-agent/util/pathlib"
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
	URL   string   `yaml:"url"`   // URL
}

func (loader *Loader) Load(path pathlib.Path, callback func(node *graph.Node) error) error {
	slog.Debug(
		"Loading node file",
		slog.String("loader", "hybrid"),
		slog.String("path", path.String()),
	)

	multiFile, err := multifile.Load(path)
	if err != nil {
		return err
	}

	metadataSegment := multiFile.FindSegmentOfType("metadata")
	if metadataSegment == nil {
		return fmt.Errorf("missing metadata segment in %s", path.String())
	}

	metadata, err := parseMetadata(metadataSegment.Contents)
	if err != nil {
		return fmt.Errorf("failed to parse metadata from %s: %v", path.String(), err)
	}

	node := graph.Node{
		Name:      metadata.Name,
		Type:      "hybrid",
		Links:     metadata.Links,
		Backlinks: nil,
		Info: &Info{
			Path: path,
		},
	}

	if err := callback(&node); err != nil {
		return err
	}

	return nil
}

func parseMetadata(lines []string) (metadata, error) {
	unparsedMetadata := strings.Join(lines, "\n")

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
