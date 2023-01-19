package main

import (
	"log"

	"01.alem.school/git/Taimas/forum/config"
	"01.alem.school/git/Taimas/forum/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
