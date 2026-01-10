package application

import "pkb-agent/tui"

type mode interface {
	Render() tui.Grid
	Handle(message tui.Message)
}
