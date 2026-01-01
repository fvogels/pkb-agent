package metaloader

import (
	"errors"
	"fmt"
	"log/slog"
	"pkb-agent/graph/node"
	"pkb-agent/graph/node/atom"
	"pkb-agent/graph/node/backblaze"
	"pkb-agent/graph/node/bookmark"
	"pkb-agent/graph/node/hybrid"
	pathlib "pkb-agent/util/pathlib"

	"gopkg.in/yaml.v3"
)

type Loader struct {
	loaders map[string]node.Loader
}

type entry struct {
	Loader string
	Path   string
}

func New() node.Loader {
	loaders := make(map[string]node.Loader)

	loaders["atom"] = atom.NewLoader()
	loaders["bookmark"] = bookmark.NewLoader()
	loaders["backblaze"] = backblaze.NewLoader()
	loaders["hybrid"] = hybrid.NewLoader()

	loader := Loader{
		loaders: loaders,
	}

	loaders["meta"] = &loader

	return &loader
}

func (loader *Loader) Load(path pathlib.Path, callback func(node node.RawNode) error) error {
	slog.Debug(
		"Loading node file",
		slog.String("loader", "meta"),
		slog.String("path", path.String()),
	)

	parentDirectory := path.Parent()
	source, err := path.ReadFile()
	if err != nil {
		return err
	}

	entries, err := loader.loadEntries(source)
	if err != nil {
		return err
	}

	errs := []error{}
	for _, entry := range entries {
		slog.Debug("Metaloader is processing entry",
			slog.String("path", entry.Path),
		)
		pathPattern := pathlib.New(entry.Path)

		paths, err := parentDirectory.Join(pathPattern).Glob()
		if err != nil {
			errs = append(errs, err)
		}

		for _, targetPath := range paths {
			subloaderName := entry.Loader
			subloader, err := loader.findLoader(subloaderName)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			if err := subloader.Load(targetPath, callback); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errors.Join(errs...)
}

func (loader *Loader) loadEntries(source []byte) ([]entry, error) {
	var entries []entry
	if err := yaml.Unmarshal(source, &entries); err != nil {
		return nil, fmt.Errorf("failed to load entries: %w", err)
	}

	return entries, nil
}

func (loader *Loader) findLoader(name string) (node.Loader, error) {
	ldr, ok := loader.loaders[name]
	if !ok {
		return nil, &ErrUnknownLoader{unknownLoaderName: name}
	}

	return ldr, nil
}
