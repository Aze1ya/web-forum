package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"01.alem.school/git/Taimas/forum/config"
)

type Server struct {
	srv http.Server
}

func (s *Server) Start(cfg *config.Config, handlers http.Handler) error {
	s.srv = http.Server{
		Addr:         cfg.Port,
		Handler:      handlers,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
	}
	fmt.Printf("Server starting http://%s%s\n", cfg.Host, cfg.Port)
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
