package snippet

import (
	"pkb-agent/graph"
	pathlib "pkb-agent/util/pathlib"
	"strings"

	"gopkg.in/yaml.v2"
)

type Loader struct{}

func NewLoader() *Loader {
	return &Loader{}
}

type metadata struct {
	Name       string
	Identifier string
	Links      []string
	Path       pathlib.Path
}

func (loader *Loader) Load(path pathlib.Path, callback func(node *graph.Node) error) error {
	unparsedMetadata, err := loader.extractMetadata(path)
	if err != nil {
		return err
	}

	metadata, err := loader.parseMetadata(unparsedMetadata)
	if err != nil {
		return err
	}

	node := graph.Node{
		Name:      metadata.Name,
		Links:     metadata.Links,
		Backlinks: nil,
		Extra: Extra{
			Path: path,
		},
	}

	if err := callback(&node); err != nil {
		return err
	}

	return nil
}

func (loader *Loader) extractMetadata(path pathlib.Path) (string, error) {
	source, err := path.ReadFile()
	if err != nil {
		return "", err
	}

	metadataLines := []string{}
	lineGenerator := strings.Lines(string(source))
	foundMetadataDelimiter := false
	lineGenerator(func(line string) bool {
		if strings.TrimSpace(line) == "---" {
			foundMetadataDelimiter = true
			return false
		}

		metadataLines = append(metadataLines, line)
		return true
	})

	if !foundMetadataDelimiter {
		return "", &ErrMissingSnippet{}
	}

	metadata := strings.Join(metadataLines, "")

	return metadata, nil
}

func (loader *Loader) parseMetadata(unparsedMetadata string) (metadata, error) {
	var result metadata

	if err := yaml.Unmarshal([]byte(unparsedMetadata), &result); err != nil {
		return metadata{}, err
	}

	return result, nil
}
