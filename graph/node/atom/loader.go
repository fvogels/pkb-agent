package atom

import (
	"errors"
	"fmt"
	"log/slog"
	"pkb-agent/graph/loaders"
	"pkb-agent/graph/node"
	"pkb-agent/util/pathlib"
	"pkb-agent/util/schema"

	"gopkg.in/yaml.v3"
)

type Loader struct{}

func init() {
	loaders.RegisterLoader("atom", NewLoader())
}

func NewLoader() *Loader {
	return &Loader{}
}

type configuration struct {
	path pathlib.Path
}

type entry struct {
	Name  string   `yaml:"name"`
	Links []string `yaml:"links"`
}

func (loader *Loader) Load(parentDirectory pathlib.Path, rawConfiguration any, callback func(node node.RawNode) error) error {
	configuration, err := loader.parseConfiguration(parentDirectory, rawConfiguration)
	if err != nil {
		return err
	}

	path := configuration.path

	slog.Debug(
		"Loading atoms node file",
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

	node := RawNode{
		name:  entry.Name,
		links: entry.Links,
	}

	return &node, nil
}

func (loader *Loader) parseConfiguration(parentDirectory pathlib.Path, rawConfiguration any) (*configuration, error) {
	var path string
	errs := []error{}

	schema.BindMapEntry(rawConfiguration, "path", &path, &errs)

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	result := configuration{
		path: parentDirectory.Join(pathlib.New(path)),
	}

	return &result, nil
}
