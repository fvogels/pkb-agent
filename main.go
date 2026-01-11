package main

import (
	"pkb-agent/tui/application"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("error loading .env file")
	}

	// cli.RunCLI()

	verbose := true
	app := application.NewApplication(verbose)
	defer app.Close()
	app.Start()
}
