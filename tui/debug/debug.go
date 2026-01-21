package debug

import (
	"log/slog"
	"pkb-agent/tui"
	"runtime"
)

func LogMessage(message tui.Message) {
	_, file, line, _ := runtime.Caller(1)
	slog.Debug(
		"Handling message",
		slog.String("file", file),
		slog.Int("line", line),
		slog.String("message", message.String()),
	)
}
