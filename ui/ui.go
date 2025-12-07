package ui

import (
	"fmt"
	"log/slog"
	"os"
	"pkb-agent/ui/screens/mainscreen"

	tea "github.com/charmbracelet/bubbletea"
)

func Start() error {
	logFile, err := os.Create("ui.log")
	if err != nil {
		fmt.Println("Failed to create log")
	}
	defer logFile.Close()

	logger := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	model := mainscreen.New()

	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		return err
	}

	return nil
}
