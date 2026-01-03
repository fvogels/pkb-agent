package tagger

import (
	"errors"
	"fmt"
	"log/slog"
	"pkb-agent/graph/loaders"
	"pkb-agent/graph/node"
	pathlib "pkb-agent/util/pathlib"
	"pkb-agent/util/schema"
)

type Loader struct{}

type configuration struct {
	extraLinks []string
	loader     string
	arguments  any
}

func init() {
	loaders.RegisterLoader("tagger", New())
}

func New() node.Loader {
	return &Loader{}
}

func (loader *Loader) Load(parentDirectory pathlib.Path, rawConfiguration any, callback func(node node.RawNode) error) error {
	configuration, err := loader.parseConfiguration(rawConfiguration)
	if err != nil {
		return fmt.Errorf("failure in tagger loader: %w", err)
	}

	subloaderName := configuration.loader
	subloaderConfiguration := configuration.arguments
	extraLinks := configuration.extraLinks

	slog.Debug(
		"Loading node file",
		slog.String("loader", "tagger"),
		slog.String("parentDirectory", parentDirectory.String()),
	)

	subloader, err := loaders.GetLoader(subloaderName)
	if err != nil {
		return fmt.Errorf("failed in tagger loader: %w", err)
	}

	wrappedCallback := func(node node.RawNode) error {
		wrapper := NodeWrapper{
			RawNode:    node,
			extraLinks: extraLinks,
		}

		return callback(&wrapper)
	}

	if err := subloader.Load(parentDirectory, subloaderConfiguration, wrappedCallback); err != nil {
		return fmt.Errorf("failed in tagger loader: %w", err)
	}

	return nil
}

func (loader *Loader) parseConfiguration(rawConfiguration any) (*configuration, error) {
	errs := []error{}
	var rawExtraLinks []any
	var result configuration

	slog.Debug("tagger", "config", fmt.Sprintf("%v", rawConfiguration))

	schema.BindMapEntry(rawConfiguration, "extraLinks", &rawExtraLinks, &errs)
	if len(errs) == 0 {
		schema.BindSlice(rawExtraLinks, &result.extraLinks, &errs)
	}
	schema.BindMapEntry(rawConfiguration, "loader", &result.loader, &errs)
	schema.BindMapEntry(rawConfiguration, "arguments", &result.arguments, &errs)

	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to read configuration: %w", errors.Join(errs...))
	}

	return &result, nil
}
