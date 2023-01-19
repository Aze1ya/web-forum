package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"01.alem.school/git/Taimas/forum/config"
	v1 "01.alem.school/git/Taimas/forum/internal/controller/http/v1"
	"01.alem.school/git/Taimas/forum/internal/repo"
	"01.alem.school/git/Taimas/forum/internal/usecase"
	"01.alem.school/git/Taimas/forum/pkg/httpserver"
	"01.alem.school/git/Taimas/forum/pkg/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func Run(cfg *config.Config) {
	db, err := sqlite.New(cfg.DbFile)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - sqlite.New: %w", err))
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("can't close db err: %v\n", err)
		} else {
			log.Printf("db closed")
		}
	}()

	repository := repo.NewRepository(db)
	useCases := usecase.NewUseCases(repository)
	handler := v1.NewHandler(useCases)
	handler.RegisterHTTPEndpoints()
	httpserver := new(httpserver.Server)

	go func() {
		if err := httpserver.Start(cfg, handler.Mux); err != nil {
			log.Println(err)
		}
	}()
	fmt.Println("Run succeed")

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	if err = httpserver.Shutdown(ctx); err != nil {
		log.Println(err)
		return
	}
}
