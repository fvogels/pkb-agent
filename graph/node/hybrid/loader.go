package hybrid

import (
	"fmt"
	"log/slog"
	"pkb-agent/graph/node"
	"pkb-agent/util/multifile"
	"pkb-agent/util/pathlib"
	"strings"

	"github.com/stretchr/testify/assert/yaml"
)

type Loader struct{}

func NewLoader() *Loader {
	return &Loader{}
}

type metadata struct {
	Name    string              `yaml:"name"`    // Name of the snippet
	Links   []string            `yaml:"links"`   // Links to other nodes
	Actions []map[string]string `yaml:"actions"` // Actions that can be performed on the node
}

func (loader *Loader) Load(path pathlib.Path, callback func(node node.RawNode) error) error {
	slog.Debug(
		"Loading hybrid node file",
		slog.String("loader", "bookmark"),
		slog.String("path", path.String()),
	)

	file, err := multifile.Load(path)
	if err != nil {
		return err
	}

	metadataSegment := file.FindSegmentOfType("metadata")
	if metadataSegment == nil {
		return fmt.Errorf("%v, file %s", ErrMissingMetadata, path.String())
	}

	metadata, err := parseMetadata(metadataSegment.Contents)
	if err != nil {
		return fmt.Errorf("failed to parse metadata from %s: %w", path.String(), err)
	}

	node := RawNode{
		name:  metadata.Name,
		links: metadata.Links,
		path:  path,
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
		slog.Error("Hybrid node is missing name", slog.String("metadata", unparsedMetadata))
		return metadata{}, ErrMissingName
	}

	return result, nil
}
