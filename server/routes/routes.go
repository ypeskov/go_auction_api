package routes

import (
	"ypeskov/go_hillel_9/internal/config"
	"ypeskov/go_hillel_9/internal/database"
	"ypeskov/go_hillel_9/internal/log"
	"ypeskov/go_hillel_9/repository/repositories"
	"ypeskov/go_hillel_9/services"
)

type Routes struct {
	Log          *log.Logger
	cfg          *config.Config
	UsersService services.UsersServiceInterface
	ItemsService services.ItemsServiceInterface
}

func New(log *log.Logger, db database.Database, cfg *config.Config) *Routes {
	itemsRepo := repositories.GetItemRepository(log, db)
	userRepo := repositories.GetUserRepository(log, db)

	return &Routes{
		Log:          log,
		cfg:          cfg,
		ItemsService: services.GetItemService(itemsRepo, log, cfg),
		UsersService: services.GetUserService(userRepo, log, cfg),
	}
}
