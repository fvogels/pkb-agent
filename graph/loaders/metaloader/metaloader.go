package metaloader

import (
	"fmt"
	"pkb-agent/graph"
	pathlib "pkb-agent/util/pathlib"

	"gopkg.in/yaml.v2"
)

type Loader struct {
	loaders map[string]graph.Loader
}

type entry struct {
	Loader string
	Path   string
}

func New() *Loader {
	return &Loader{
		loaders: make(map[string]graph.Loader),
	}
}

func (loader *Loader) AddLoader(name string, ldr graph.Loader) error {
	_, ok := loader.loaders[name]
	if ok {
		return &ErrDuplicateLoaderName{duplicateLoaderName: name}
	}

	loader.loaders[name] = ldr
	return nil
}

func (loader *Loader) Load(path pathlib.Path, callback func(node graph.Node) error) error {
	parentDirectory := path.Parent()
	source, err := path.ReadFile()
	if err != nil {
		return err
	}

	entries, err := loader.loadEntries(source)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		pathPattern := pathlib.New(entry.Path)

		paths, err := parentDirectory.Join(pathPattern).Glob()
		if err != nil {
			return fmt.Errorf("error globbing %s: %w", pathPattern, err)
		}

		for _, targetPath := range paths {
			subloaderName := entry.Loader
			subloader, err := loader.findLoader(subloaderName)
			if err != nil {
				return err
			}

			if err := subloader.Load(targetPath, callback); err != nil {
				return err
			}
		}
	}

	return nil
}

func (loader *Loader) loadEntries(source []byte) ([]entry, error) {
	var entries []entry
	if err := yaml.Unmarshal(source, &entries); err != nil {
		return nil, fmt.Errorf("failed to load entries: %w", err)
	}

	return entries, nil
}

func (loader *Loader) findLoader(name string) (graph.Loader, error) {
	ldr, ok := loader.loaders[name]
	if !ok {
		return nil, &ErrUnknownLoader{unknownLoaderName: name}
	}

	return ldr, nil
}
