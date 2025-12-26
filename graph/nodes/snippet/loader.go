package snippet

import (
	"fmt"
	"log/slog"
	"pkb-agent/graph"
	pathlib "pkb-agent/util/pathlib"
	"strings"

	"gopkg.in/yaml.v3"
)

type Loader struct{}

func NewLoader() *Loader {
	return &Loader{}
}

type metadata struct {
	Name              string   `yaml:"name"`      // Name of the snippet.
	Language          string   `yaml:"language"`  // Language of the snippet.
	HighlightLanguage string   `yaml:"highlight"` // Optional: language identifier used for syntax highlighting. If missing, a lowercase version of language will be used.
	Links             []string `yaml:"links"`     // Links to other nodes.
}

func (loader *Loader) Load(path pathlib.Path, callback func(node *graph.Node) error) error {
	slog.Debug(
		"Loading node file",
		slog.String("loader", "snippet"),
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

	if len(metadata.Name) == 0 {
		slog.Error("Snippet node is missing name", slog.String("path", path.String()))
		return &ErrMissingName{path: path}
	}

	var highlightLanguage string
	if len(metadata.HighlightLanguage) == 0 {
		highlightLanguage = strings.ToLower(metadata.Language)
	} else {
		highlightLanguage = metadata.HighlightLanguage
	}

	node := graph.Node{
		Name:      metadata.Name,
		Type:      "snippet",
		Links:     append(metadata.Links, "Snippet"),
		Backlinks: nil,
		Info: &Info{
			Path:                    path,
			LanguageForHighlighting: highlightLanguage,
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
		return "", &ErrMissingSnippet{path: path}
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
