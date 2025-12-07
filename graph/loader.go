package graph

import pathlib "pkb-agent/util/path"

type Loader interface {
	Load(path pathlib.Path, callback func(node Node) error) error
}
