package server

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "ypeskov/go_hillel_9/docs"
	"ypeskov/go_hillel_9/internal/config"
	"ypeskov/go_hillel_9/internal/log"
	"ypeskov/go_hillel_9/server/middleware"
	"ypeskov/go_hillel_9/server/routes"
)

type Server struct {
	e    *echo.Echo
	port string
	log  *log.Logger
}

func New(cfg *config.Config, handlers *routes.Routes) *Server {
	e := echo.New()

	// e.Use(middleware.Logger())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	itemsGroup := e.Group("/items")
	itemsGroup.Use(middleware.AuthMiddleware(handlers.Log, cfg, handlers.UsersService))
	handlers.RegisterItemsRoutes(itemsGroup)

	usersGroup := e.Group("/users")
	handlers.RegisterUsersRoutes(usersGroup)

	return &Server{
		e:    e,
		port: cfg.Port,
		log:  handlers.Log,
	}
}

func (s *Server) Start() error {
	s.log.Infof("Starting the server on port %s", s.port)

	return s.e.Start(s.port)
}
