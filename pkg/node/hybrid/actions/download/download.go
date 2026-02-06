package download

import (
	"fmt"
	"pkb-agent/extern"
	"pkb-agent/pkg/node"
)

type Action struct {
	description string
}

type B2 struct {
	Action
	bucket   string
	filename string
}

func (action Action) GetDescription() string {
	return action.description
}

func (action B2) Perform() {
	extern.BackblazeDownloadAndOpen(action.bucket, action.filename)
}

func Parse(rawAction map[string]string) (node.Action, error) {
	description, err := getAttribute(rawAction, "description")
	if err != nil {
		return nil, err
	}

	source, err := getAttribute(rawAction, "source")
	if err != nil {
		return nil, err
	}

	switch source {
	case "backblaze":
		bucket, err := getAttribute(rawAction, "bucket")
		if err != nil {
			return nil, err
		}

		filename, err := getAttribute(rawAction, "filename")
		if err != nil {
			return nil, err
		}

		return &B2{
			Action: Action{
				description: description,
			},
			bucket:   bucket,
			filename: filename,
		}, nil

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownSource, source)
	}
}

func getAttribute(rawAction map[string]string, fieldName string) (string, error) {
	value, ok := rawAction[fieldName]

	if !ok {
		return "", fmt.Errorf("%w: %s", ErrMissingAttribute, fieldName)
	}

	return value, nil
}
