package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Do I want to have the interface here or move it to the pacakge log?
type logger interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}

type silentLog struct{}

func (l silentLog) Debugf(template string, args ...interface{}) {}
func (l silentLog) Infof(template string, args ...interface{})  {}
func (l silentLog) Warnf(template string, args ...interface{})  {}
func (l silentLog) Errorf(template string, args ...interface{}) {}

type Server struct {
	server http.Server
	logger logger
}

// NewServer creates a new sever instance that can be started individually
func NewServer(addr string, log logger) *Server {

	if addr == "" {
		addr = ":8080"
	}

	l := silentLog{}
	if log == nil {
		log = &l
	}

	return &Server{
		logger: log,
		server: http.Server{
			Addr: addr,
			//Handler: apiHandler,
		},
	}
}

// StopOnOsSignal creates a channel and a routine that will stop the server when a
// SIGTERM or SIGINT is received
func (srv *Server) stopOnOsSignal() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalCh
		srv.Stop()
	}()
}

// Start to listen on the configured address
func (srv *Server) Start() error {

	srv.stopOnOsSignal()

	srv.logger.Infof("Starting server on: %s", srv.server.Addr)
	if err := srv.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (srv *Server) Stop() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.server.Shutdown(ctx); err != nil {
		srv.logger.Infof("shutdown: %v", err)
	}

}
