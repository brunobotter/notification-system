package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/brunobotter/notification-system/infra/logger"
	"github.com/brunobotter/notification-system/main/config"
	"github.com/brunobotter/notification-system/main/container"
	"github.com/brunobotter/notification-system/main/server/router"

	"github.com/labstack/echo/v4"
)

type Server struct {
	container container.Container
	config    *config.Config
	logger    logger.Logger
	echo      *echo.Echo
}

func NewServer(container container.Container) (*Server, error) {
	server := &Server{
		container: container,
		echo:      echo.New(),
	}
	container.Resolve(&server.config)

	container.Resolve(&server.logger)

	server.setup()
	return server, nil
}

func (s *Server) setup() {
	s.echo.HideBanner = true

	var cfg *config.Config
	var logger *logger.Logger
	s.container.Resolve(&cfg)
	s.container.Resolve(&logger)

	router.RegisterRouter(s.echo, cfg, s.container)
}

func (s *Server) Run(ctx context.Context) {
	go func() {
		address := fmt.Sprintf(":%d", s.config.Server.Port)
		if err := s.echo.Start(address); err != nil && err != http.ErrServerClosed {
			s.echo.Logger.Fatal(err)
		}
	}()
	s.waitForShutdown(ctx)
}

func (s *Server) waitForShutdown(ctx context.Context) {
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.echo.Shutdown(ctx); err != nil {
		s.echo.Logger.Fatal(err)
	}
}
