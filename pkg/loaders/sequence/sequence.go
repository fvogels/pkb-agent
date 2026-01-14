package sequence

import (
	"errors"
	"fmt"
	"log/slog"
	"pkb-agent/pkg/loaders"
	"pkb-agent/pkg/node"
	pathlib "pkb-agent/util/pathlib"
	"pkb-agent/util/schema"
)

type Loader struct{}

type configuration struct {
	entries []entry
}

type entry struct {
	loader    string
	arguments any
}

func init() {
	loaders.RegisterLoader("sequence", New())
}

func New() node.Loader {
	return &Loader{}
}

func (loader *Loader) Load(parentDirectory pathlib.Path, rawConfiguration any, callback func(node node.RawNode) error) error {
	configuration, err := loader.parseConfiguration(rawConfiguration)
	if err != nil {
		return fmt.Errorf("failed to load sequence: %w", err)
	}

	entries := configuration.entries

	slog.Debug(
		"Loading node file",
		slog.String("loader", "sequence"),
		slog.String("parentDirectory", parentDirectory.String()),
	)

	errs := []error{}
	for entryIndex, entry := range entries {
		subloader, err := loaders.GetLoader(entry.loader)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed in sequence loader at entry with index %d: %w", entryIndex, err))
			continue
		}

		if err := subloader.Load(parentDirectory, entry.arguments, callback); err != nil {
			errs = append(errs, fmt.Errorf("failed in sequence loader at entry with index %d: %w", entryIndex, err))
		}
	}

	return errors.Join(errs...)
}

func (loader *Loader) parseConfiguration(rawConfiguration any) (*configuration, error) {
	var items []map[string]any
	errs := []error{}

	schema.BindSlice(rawConfiguration, &items, &errs)

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	entries := make([]entry, len(items))

	for index, item := range items {
		schema.BindMapEntry(item, "loader", &entries[index].loader, &errs)
		schema.BindMapEntry(item, "arguments", &entries[index].arguments, &errs)
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to parse configuration: %w", errors.Join(errs...))
	}

	result := configuration{
		entries: entries,
	}

	return &result, nil
}
