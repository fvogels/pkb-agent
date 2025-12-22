package hybrid

import "fmt"

type DownloadAction struct {
	ActionBase
	Bucket   string
	Filename string
}

func parseDownloadAction(metadata map[string]string) (*DownloadAction, error) {
	description, found := metadata["description"]
	if !found {
		return nil, fmt.Errorf("missing description: %w", &ErrInvalidAction{})
	}

	source, found := metadata["source"]
	if !found {
		return nil, fmt.Errorf("missing source: %w", &ErrInvalidAction{})
	}
	if source != "backblaze" {
		return nil, fmt.Errorf("source should be backblaze: %w", &ErrInvalidAction{})
	}

	bucket, found := metadata["bucket"]
	if !found {
		return nil, fmt.Errorf("missing bucket: %w", &ErrInvalidAction{})
	}

	filename, found := metadata["filename"]
	if !found {
		return nil, fmt.Errorf("missing filename: %w", &ErrInvalidAction{})
	}

	action := DownloadAction{
		ActionBase: ActionBase{
			Description: description,
		},
		Bucket:   bucket,
		Filename: filename,
	}

	return &action, nil
}
