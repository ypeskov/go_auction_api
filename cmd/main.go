package main

import (
	"fmt"
	"ypeskov/go_hillel_9/internal/config"
	"ypeskov/go_hillel_9/internal/database"
	log "ypeskov/go_hillel_9/internal/log"
	"ypeskov/go_hillel_9/server"
	"ypeskov/go_hillel_9/server/routes"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)

		return
	}

	logger := log.New(cfg)
	logger.Info("Starting the application...")

	db := database.GetDB(cfg, logger)

	routes := routes.New(logger, db, cfg)

	server := server.New(cfg, routes)
	err = server.Start()
	if err != nil {
		logger.Errorf("Error starting the server: %v", err)
	}
}
