package routes

import (
	"ypeskov/go_hillel_9/internal/database"
	"ypeskov/go_hillel_9/internal/log"
	"ypeskov/go_hillel_9/repository/repositories"
	"ypeskov/go_hillel_9/services"
)

type Routes struct {
	Log          *log.Logger
	ItemsRepo    repositories.ItemRepositoryInterface
	UserRepo     repositories.UserRepositoryInterface
	usersService services.UsersServiceInterface
}

func New(log *log.Logger, db database.Database) *Routes {
	itemsRepo := repositories.GetItemRepository(log, db)
	userRepo := repositories.GetUserRepository(log, db)

	return &Routes{
		Log:          log,
		ItemsRepo:    itemsRepo,
		UserRepo:     userRepo,
		usersService: services.GetUserService(userRepo, log),
	}
}
