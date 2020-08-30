package internalhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/avbru/homeworks/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	srv   *http.Server
	app   Application
	shutT time.Duration
}

type Application interface {
	// TODO
}

func NewServer(app Application, addr string, shutT time.Duration) *Server {
	return &Server{
		app:   app,
		shutT: shutT,
		srv: &http.Server{
			Addr:    addr,
			Handler: router(),
		},
	}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutT)
	cancel()

	return s.srv.Shutdown(ctx)
}

func router() http.Handler {
	h := http.NewServeMux()
	h.Handle("/hello", loggingMiddleware(http.HandlerFunc(helloHandler)))

	return h
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("hello world"))
	if err != nil {
		logger.Logger.Err(err)
	}
}
