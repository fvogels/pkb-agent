package metaloader

import (
	"errors"
	"log/slog"
	"pkb-agent/pkg/node"
	"pkb-agent/pkg/node/atom"
	"pkb-agent/pkg/node/backblaze"
	"pkb-agent/pkg/node/bookmark"
	"pkb-agent/pkg/node/hybrid"
	pathlib "pkb-agent/util/pathlib"
	"pkb-agent/util/schema"
)

type Loader struct {
	loaders map[string]node.Loader
}

type entry struct {
	Loader    string
	Arguments any
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

func (loader *Loader) Load(parentDirectory pathlib.Path, configuration any, callback func(node node.RawNode) error) error {
	slog.Debug(
		"Loading node file",
		slog.String("loader", "meta"),
	)

	entries, err := loader.parseConfiguration(configuration)
	if err != nil {
		return err
	}

	errs := []error{}
	for _, entry := range entries {
		subloader, err := loader.findLoader(entry.Loader)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if err := subloader.Load(parentDirectory, entry.Arguments, callback); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (loader *Loader) parseConfiguration(configuration any) ([]entry, error) {
	var items []map[string]any
	errs := []error{}

	schema.BindSlice(configuration, &items, &errs)
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	entries := make([]entry, len(items))
	for index, item := range items {
		schema.BindMapEntry(item, "loader", &entries[index].Loader, &errs)
		schema.BindMapEntry(item, "arguments", &entries[index].Arguments, &errs)
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
