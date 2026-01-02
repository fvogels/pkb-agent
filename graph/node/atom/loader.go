package atom

import (
	"errors"
	"fmt"
	"log/slog"
	"pkb-agent/graph/node"
	"pkb-agent/util/pathlib"

	"gopkg.in/yaml.v3"
)

type Loader struct{}

func NewLoader() *Loader {
	return &Loader{}
}

type entry struct {
	Name  string   `yaml:"name"`
	Links []string `yaml:"links"`
}

func (loader *Loader) Load(path pathlib.Path, callback func(node node.RawNode) error) error {
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
