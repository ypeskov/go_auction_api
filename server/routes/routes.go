package routes

import (
	"ypeskov/go_hillel_9/internal/database"
	"ypeskov/go_hillel_9/internal/log"
	"ypeskov/go_hillel_9/repository/repositories"
)

type Routes struct {
	Log       *log.Logger
	ItemsRepo *repositories.ItemRepository
}

func New(log *log.Logger, db database.Database) *Routes {
	itemsRepo := repositories.NewItemRepository(log, db)

	return &Routes{
		Log:       log,
		ItemsRepo: itemsRepo,
	}
}
