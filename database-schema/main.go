package main

import (
	"log"

	"01.alem.school/git/Taimas/forum/config"
	"01.alem.school/git/Taimas/forum/database-schema/up"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	err = up.DbSqliteUp(cfg)
	if err != nil {
		log.Fatalf("DB up error: %s", err)
	}
}
