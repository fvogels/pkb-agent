package backblaze

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

type bucket struct {
	Name  string  `yaml:"bucket"`
	Files []entry `yaml:"files"`
}

type entry struct {
	Name     string
	Links    []string
	Filename string
}

type configuration struct {
	path pathlib.Path
}

func init() {
	loaders.RegisterLoader("backblaze", NewLoader())
}

func NewLoader() *Loader {
	return &Loader{}
}

func (loader *Loader) Load(parentDirectory pathlib.Path, rawConfiguration any, callback func(node node.RawNode) error) error {
	configuration, err := loader.parseConfiguration(parentDirectory, rawConfiguration)
	path := configuration.path

	slog.Debug(
		"Loading node file",
		slog.String("loader", "backblaze"),
		slog.String("path", path.String()),
	)

	source, err := path.ReadFile()
	if err != nil {
		return err
	}

	var buckets []bucket
	if err := yaml.Unmarshal(source, &buckets); err != nil {
		return err
	}

	var errs []error
	for _, bucket := range buckets {
		if len(bucket.Name) == 0 {
			errs = append(errs, ErrBucketMissing)
		}

		for fileIndex, file := range bucket.Files {
			if len(file.Name) == 0 {
				err := fmt.Errorf("%w, path: %s, index: %d, bucket: %s", ErrNameMissing, path.String(), fileIndex, bucket.Name)
				errs = append(errs, err)
				continue
			}

			if len(file.Filename) == 0 {
				err := fmt.Errorf("%w, path: %s, index: %d, node: %s, bucket: %s", ErrFilenameMissing, path.String(), fileIndex, file.Name, bucket.Name)
				errs = append(errs, err)
				continue
			}

			node := RawNode{
				name:     file.Name,
				bucket:   bucket.Name,
				filename: file.Filename,
				links:    file.Links,
			}

			if err := callback(&node); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errors.Join(errs...)
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
