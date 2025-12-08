package pathlib

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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
	absolute, err := p.Absolute()
	if err != nil {
		return nil, err
	}

	if prefix, ok := strings.CutSuffix(absolute.path, "*"); ok {
		result := []Path{}
		walker := func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.Mode().IsRegular() {
				result = append(result, New(path))
			}

			return nil
		}

		filepath.Walk(prefix, walker)
		return result, nil
	} else if strings.Contains(p.path, "*") {
		return nil, fmt.Errorf("unsupported glob pattern")
	} else {
		return []Path{p}, nil
	}
}

func (p Path) ReadFile() ([]byte, error) {
	return os.ReadFile(p.path)
}

func (p Path) String() string {
	return p.path
}
