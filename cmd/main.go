package main

import (
	"ypeskov/go_hillel_9/server"
	"ypeskov/go_hillel_9/server/routes"

	"fmt"
	"ypeskov/go_hillel_9/internal/config"
	"ypeskov/go_hillel_9/internal/logger"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	logger := logger.New(cfg)
	logger.Info("Starting the application...")

	routes := routes.New(logger)

	server := server.New(cfg, routes)
	err = server.Start()
	if err != nil {
		logger.Errorf("Error starting the server: %v", err)
	}
}
