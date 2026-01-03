package hybrid

import (
	"errors"
	"fmt"
	"log/slog"
	"pkb-agent/graph/loaders"
	"pkb-agent/graph/node"
	"pkb-agent/util/multifile"
	"pkb-agent/util/pathlib"
	"pkb-agent/util/schema"
	"strings"

	"github.com/stretchr/testify/assert/yaml"
)

type Loader struct{}

type metadata struct {
	Name    string              `yaml:"name"`    // Name of the snippet
	Links   []string            `yaml:"links"`   // Links to other nodes
	Actions []map[string]string `yaml:"actions"` // Actions that can be performed on the node
}

type configuration struct {
	path pathlib.Path
}

func init() {
	loaders.RegisterLoader("hybrid", NewLoader())
}

func NewLoader() *Loader {
	return &Loader{}
}

func (loader *Loader) Load(parentDirectory pathlib.Path, rawConfiguration any, callback func(node node.RawNode) error) error {
	configuration, err := loader.parseConfiguration(parentDirectory, rawConfiguration)
	if err != nil {
		return fmt.Errorf("failed to load hybrid node: %w", err)
	}

	path := configuration.path

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

func (loader *Loader) parseConfiguration(parentDirectory pathlib.Path, rawConfiguration any) (*configuration, error) {
	var path string
	var errs = []error{}

	schema.BindMapEntry(rawConfiguration, "path", &path, &errs)

	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to read configuration: %w", errors.Join(errs...))
	}

	configuration := configuration{
		path: parentDirectory.Join(pathlib.New(path)),
	}

	return &configuration, nil
}
