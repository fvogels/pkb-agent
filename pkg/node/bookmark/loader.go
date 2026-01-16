package bookmark

import (
	"errors"
	"fmt"
	"log/slog"
	"pkb-agent/pkg/loaders"
	"pkb-agent/pkg/node"
	"pkb-agent/util/pathlib"
	"pkb-agent/util/schema"

	"gopkg.in/yaml.v3"
)

type Loader struct{}

type entry struct {
	Name        string   `yaml:"name"`
	Links       []string `yaml:"links"`
	URL         string   `yaml:"url"`
	Description string   `yaml:"description"`
}

type configuration struct {
	path pathlib.Path
}

func init() {
	loaders.RegisterLoader("bookmark", NewLoader())
}

func NewLoader() *Loader {
	return &Loader{}
}

func (loader *Loader) Load(parentDirectory pathlib.Path, rawConfiguration any, callback func(node node.RawNode) error) error {
	configuration, err := loader.parseConfiguration(parentDirectory, rawConfiguration)
	path := configuration.path

	slog.Debug(
		"Loading bookmarks node file",
		slog.String("loader", "bookmark"),
		slog.String("path", path.String()),
	)

	fileContents, err := loader.readFile(path)
	if err != nil {
		return err
	}

	entries, err := loader.parseFileContents(fileContents)
	if err != nil {
		return err
	}

	var errs []error
	for index, entry := range entries {
		rawNode, err := loader.convertEntryToNode(&entry)

		if err != nil {
			err = fmt.Errorf("File: %s, index: %d, error: %w", path.String(), index, err)
			errs = append(errs, err)
			continue
		}

		if err := callback(rawNode); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (loader *Loader) readFile(path pathlib.Path) ([]byte, error) {
	return path.ReadFile()
}

func (loader *Loader) parseFileContents(data []byte) ([]entry, error) {
	var entries []entry

	if err := yaml.Unmarshal(data, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

func (loader *Loader) convertEntryToNode(entry *entry) (*RawNode, error) {
	if len(entry.Name) == 0 {
		return nil, ErrMissingName
	}

	if len(entry.URL) == 0 {
		return nil, ErrMissingURL
	}

	if len(entry.Description) == 0 {
		return nil, ErrMissingDescription
	}

	node := RawNode{
		name:  entry.Name,
		url:   entry.URL,
		links: entry.Links,
	}

	return &node, nil
}

func (loader *Loader) parseConfiguration(parentDirectory pathlib.Path, rawConfiguration any) (*configuration, error) {
	var path string
	var errs = []error{}

	schema.BindMapEntry(rawConfiguration, "path", &path, &errs)

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	configuration := configuration{
		path: parentDirectory.Join(pathlib.New(path)),
	}

	return &configuration, nil
}
