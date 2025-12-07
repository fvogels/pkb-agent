package pathlib

import (
	"errors"
	"os"
	"path/filepath"
	"pkb-agent/util"
)

type Path struct {
	path string
}

func New(path string) Path {
	return Path{path: path}
}

func (p Path) Join(other Path) Path {
	return New(filepath.Join(p.path, other.path))
}

func (p Path) IsDirectory() (bool, error) {
	stat, err := os.Stat(p.path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, err
	}

	return stat.Mode().IsDir(), nil
}

func (p Path) IsFile() (bool, error) {
	stat, err := os.Stat(p.path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, err
	}

	return stat.Mode().IsRegular(), nil
}

func (p Path) Absolute() (Path, error) {
	absolute, err := filepath.Abs(p.path)
	if err != nil {
		return Path{}, err
	}

	return New(absolute), nil
}

func (p Path) Parent() Path {
	return New(filepath.Dir(p.path))
}

func (p Path) Basename() string {
	return filepath.Base(p.path)
}

func (p Path) Glob() ([]Path, error) {
	paths, err := filepath.Glob(p.path)

	if err != nil {
		return nil, err
	}

	return util.Map(paths, New), nil
}

func (p Path) ReadFile() ([]byte, error) {
	return os.ReadFile(p.path)
}
