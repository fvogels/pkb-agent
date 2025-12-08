package graph

type ErrNameClash struct{}

func (err *ErrNameClash) Error() string {
	return "name clash"
}

type ErrUnknownNodes struct{}

func (err *ErrUnknownNodes) Error() string {
	return "unknown node"
}
