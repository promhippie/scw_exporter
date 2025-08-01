package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/promhippie/scw_exporter/pkg/command"
)

func main() {
	if env := os.Getenv("SCW_EXPORTER_ENV_FILE"); env != "" {
		_ = godotenv.Load(env)
	}

	if err := command.Run(); err != nil {
		os.Exit(1)
	}
}
