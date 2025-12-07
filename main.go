package main

import (
	"pkb-agent/cli"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("error loading .env file")
	}

	cli.RunCLI()
}
