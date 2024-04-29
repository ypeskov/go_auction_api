package routes

import (
	"ypeskov/go_hillel_9/internal/database"
	"ypeskov/go_hillel_9/internal/logger"
)

type Routes struct {
	Log *logger.Logger

	Db *database.Database
}

func New(log *logger.Logger, dbConnection *database.Database) *Routes {
	return &Routes{
		Log: log,
		Db:  dbConnection,
	}
}
