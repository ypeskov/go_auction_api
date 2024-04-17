package server

import (
	"fmt"
	"ypeskov/go_hillel_9/internal/config"
	"ypeskov/go_hillel_9/internal/logger"
	"ypeskov/go_hillel_9/server/routes"

	"github.com/labstack/echo/v4"
)

type Server struct {
	e    *echo.Echo
	port string
	log  *logger.Logger
}

func New(cfg *config.Config, handlers *routes.Routes) *Server {
	e := echo.New()

	itemsGroup := e.Group("/items")
	handlers.RegisterItemsRoutes(itemsGroup)

	return &Server{
		e:    e,
		port: cfg.Port,
		log:  handlers.Log,
	}
}

func (s *Server) Start() error {
	s.log.Info("Starting the server...")
	fmt.Printf("port: %s\n", s.port)
	return s.e.Start(s.port)
}
