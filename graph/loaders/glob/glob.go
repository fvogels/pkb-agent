package glob

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
	path      pathlib.Path
	loader    string
	arguments map[string]any
}

func init() {
	loaders.RegisterLoader("glob", New())
}

func New() node.Loader {
	return &Loader{}
}

func (loader *Loader) Load(parentDirectory pathlib.Path, rawConfiguration any, callback func(node node.RawNode) error) error {
	configuration, err := loader.parseConfiguration(parentDirectory, rawConfiguration)
	if err != nil {
		return fmt.Errorf("glob loader failed: %w", err)
	}

	rootPath := configuration.path
	subloaderName := configuration.loader
	subloaderConfiguration := configuration.arguments

	slog.Debug(
		"Loading node file",
		slog.String("loader", "glob"),
		slog.String("parentDirectory", parentDirectory.String()),
		slog.String("subloader", subloaderName),
		slog.String("rootPath", rootPath.String()),
	)

	subloader, err := loaders.GetLoader(subloaderName)
	if err != nil {
		return fmt.Errorf("glob loader failed to list files in %s: %w", parentDirectory.String(), err)
	}

	paths, err := rootPath.FindFiles()
	if err != nil {
		return fmt.Errorf("glob loader failed to list files in %s: %w", parentDirectory.String(), err)
	}

	errs := []error{}
	for _, path := range paths {
		subloaderConfiguration["path"] = path.Basename()
		if err := subloader.Load(path.Parent(), subloaderConfiguration, callback); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (loader *Loader) parseConfiguration(parentDirectory pathlib.Path, rawConfiguration any) (*configuration, error) {
	var path string
	var subloader string
	var subloaderArguments map[string]any
	errs := []error{}

	schema.BindMapEntry(rawConfiguration, "path", &path, &errs)
	schema.BindMapEntry(rawConfiguration, "loader", &subloader, &errs)
	schema.BindMapEntry(rawConfiguration, "arguments", &subloaderArguments, &errs)

	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to parse configuration: %w", errors.Join(errs...))
	}

	result := configuration{
		path:      parentDirectory.Join(pathlib.New(path)),
		loader:    subloader,
		arguments: subloaderArguments,
	}

	return &result, nil
}
