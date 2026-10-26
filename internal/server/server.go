package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"strconv"

	"chatsmith-mini/configs"
	"chatsmith-mini/constants"
	"chatsmith-mini/pkg/failure"
)

const (
	listenNetwork     = "tcp"
	logServerStart    = "server_start"
	logServerShutdown = "server_shutdown"
)

type Server struct {
	cfg  configs.ServerConfig
	log  *slog.Logger
	http *http.Server
}

func New(cfg configs.ServerConfig, log *slog.Logger, handler http.Handler) *Server {
	return &Server{
		cfg: cfg,
		log: log,
		http: &http.Server{
			Addr:              ":" + strconv.Itoa(cfg.Port),
			Handler:           handler,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
			ReadTimeout:       cfg.ReadTimeout,
			IdleTimeout:       cfg.IdleTimeout,
			MaxHeaderBytes:    cfg.MaxHeaderBytes,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	listener, err := net.Listen(listenNetwork, s.http.Addr)
	if err != nil {
		return failure.ErrWithTrace(err)
	}

	failed := make(chan error, 1)
	go func() { failed <- s.http.Serve(listener) }()

	s.log.Info(logServerStart, constants.LogKeyPort, s.cfg.Port)

	select {
	case err := <-failed:
		return failure.ErrWithTrace(err)
	case <-ctx.Done():
	}

	s.log.Info(logServerShutdown)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()

	if err := s.http.Shutdown(shutdownCtx); err != nil {
		return failure.ErrWithTrace(err)
	}

	return nil
}
