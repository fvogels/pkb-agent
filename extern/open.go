package extern

import (
	"fmt"
	"os/exec"
	"pkb-agent/util/pathlib"
)

func OpenUsingDefaultViewer(path pathlib.Path) error {
	absolutePath, err := path.Absolute()
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	command := exec.Command("rundll32", "url.dll,FileProtocolHandler", absolutePath.String())
	if err := command.Start(); err != nil {
		return err
	}

	return nil
}
