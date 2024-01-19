package server

import (
	"context"
	"errors"
	"fmt"
	"git.andresbott.com/Golang/carbon/app/server/handlers"
	"git.andresbott.com/Golang/carbon/libs/log/zero"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	server http.Server
	logger *zerolog.Logger
}

type Cfg struct {
	Addr   string
	Logger *zerolog.Logger
	Db     *gorm.DB
}

// NewServer creates a new sever instance that can be started individually
func NewServer(cfg Cfg) *Server {

	if cfg.Addr == "" {
		cfg.Addr = ":8080"
	}

	if cfg.Logger == nil {
		cfg.Logger = zero.Silent()
	}

	handler := handlers.NewAppHandler(cfg.Logger, cfg.Db)
	if handler == nil {
		panic("nil")
	}

	return &Server{
		logger: cfg.Logger,
		server: http.Server{
			Addr:    cfg.Addr,
			Handler: handler,
		},
	}
}

// Start to listen on the configured address
func (srv *Server) Start() error {

	done := make(chan bool, 1)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalCh
		srv.Stop()
		done <- true
	}()

	srv.logger.Info().Msg(fmt.Sprintf("Starting server on: %s", srv.server.Addr))
	if err := srv.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	<-done
	return nil
}

// Stop shut down the server cleanly
func (srv *Server) Stop() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.server.Shutdown(ctx); err != nil {
		srv.logger.Warn().Msg(fmt.Sprintf("shutdown: %v", err))
	}
	srv.logger.Info().Msg("server stopped")

}
