package server

import (
	"context"
	"errors"
	"fmt"
	"git.andresbott.com/Golang/carbon/app/server/handlers"
	"git.andresbott.com/Golang/carbon/libs/log/zero"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	server    http.Server
	obsServer *http.Server
	logger    *zerolog.Logger
}

type Cfg struct {
	Addr    string
	ObsAddr string
	Logger  *zerolog.Logger
	Db      *gorm.DB
}

// NewServer creates a new sever instance that can be started individually
func NewServer(cfg Cfg) *Server {

	if cfg.Addr == "" {
		cfg.Addr = ":8080"
	}
	if cfg.ObsAddr == "" {
		cfg.ObsAddr = ":9090"
	}

	if cfg.Logger == nil {
		cfg.Logger = zero.Silent()
	}

	handler, err := handlers.NewAppHandler(cfg.Logger, cfg.Db)
	if err != nil {
		panic(fmt.Sprintf("unable to initialize app: %v", err))
	}

	return &Server{
		logger: cfg.Logger,
		server: http.Server{
			Addr:    cfg.Addr,
			Handler: handler,
		},

		obsServer: &http.Server{
			Addr:    cfg.ObsAddr,
			Handler: obsHandler(),
		},
	}
}

// Start to listen on the configured address
func (srv *Server) Start() error {

	stopDone := make(chan bool, 1)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// handle shutdown
	go func() {
		<-signalCh
		srv.Stop()
		stopDone <- true
	}()

	// observability server
	go func() {
		ln, err := net.Listen("tcp", srv.obsServer.Addr)
		if err != nil {
			panic(fmt.Sprintf("error starting obserbability server: %v", err))
		}
		srv.logger.Info().Msg(fmt.Sprintf("obserbability server started on: %s", srv.obsServer.Addr))

		if err := srv.obsServer.Serve(ln); !errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Sprintf("error in obserbability server: %v", err))
		}

	}()

	ln, err := net.Listen("tcp", srv.server.Addr)
	if err != nil {
		return err
	}
	srv.logger.Info().Msg(fmt.Sprintf("server started on: %s", srv.server.Addr))

	if err = srv.server.Serve(ln); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	<-stopDone
	return nil
}

// Stop shut down the server cleanly
func (srv *Server) Stop() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.server.Shutdown(ctx); err != nil {
		srv.logger.Warn().Msg(fmt.Sprintf("server shutdown: %v", err))
	}
	srv.logger.Info().Msg("server stopped")

	if err := srv.obsServer.Shutdown(ctx); err != nil {
		srv.logger.Warn().Msg(fmt.Sprintf("obs server shutdown: %v", err))
	}
	srv.logger.Info().Msg("obs server stopped")
}

func obsHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		content := `
<a href="/metrics">/metrics</a>
`
		_, _ = fmt.Fprint(writer, content)

	})
	mux.Handle("/metrics", promhttp.Handler())
	return mux
}
