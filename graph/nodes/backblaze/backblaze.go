package backblaze

import (
	"errors"
	"log/slog"
	"pkb-agent/graph"
	"pkb-agent/util/pathlib"

	"gopkg.in/yaml.v3"
)

type Loader struct{}

type Extra struct {
	BucketName string
	Filename   string
}

func NewLoader() *Loader {
	return &Loader{}
}

type bucket struct {
	Name  string  `yaml:"bucket"`
	Files []entry `yaml:"files"`
}

type entry struct {
	Name     string
	Links    []string
	Filename string
}

func (loader *Loader) Load(path pathlib.Path, callback func(node *graph.Node) error) error {
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
	for bucketIndex, bucket := range buckets {
		if len(bucket.Name) == 0 {
			errs = append(errs, &ErrBucketMissingName{path: path, index: bucketIndex})
		}

		for fileIndex, file := range bucket.Files {
			if len(file.Name) == 0 {
				errs = append(errs, &ErrFileMissingName{path: path, index: fileIndex})
				continue
			}

			node := graph.Node{
				Name:      file.Name,
				Type:      "file",
				Links:     append(file.Links, "File"),
				Backlinks: nil,
				Extra: &Extra{
					BucketName: bucket.Name,
					Filename:   file.Filename,
				},
			}

			if err := callback(&node); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errors.Join(errs...)
}
