package server

import (
	"context"
	logger "git.andresbott.com/Golang/carbon/libs/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	server  http.Server
	handler *http.Handler
	logger  logger.LeveledLogger
}

// NewServer creates a new sever instance that can be started individually
func NewServer(addr string, log logger.LeveledLogger) *Server {

	if addr == "" {
		addr = ":8080"
	}

	if log == nil {
		log = &logger.SilentLog{}
	}

	handler := newMainHandler(log)
	if handler == nil {
		panic("nil")
	}

	return &Server{
		logger: log,
		server: http.Server{
			Addr:    addr,
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

	srv.logger.Infof("Starting server on: %s", srv.server.Addr)
	if err := srv.server.ListenAndServe(); err != http.ErrServerClosed {
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
		srv.logger.Infof("shutdown: %v", err)
	}
	srv.logger.Infof("server stopped %s", "")

}
