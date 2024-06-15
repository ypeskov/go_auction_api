package services

import (
	"ypeskov/go_hillel_9/internal/config"
	"ypeskov/go_hillel_9/internal/log"
	"ypeskov/go_hillel_9/repository/models"
	"ypeskov/go_hillel_9/repository/repositories"
)

type UserTypeService struct {
	log          *log.Logger
	cfg          *config.Config
	userTypeRepo repositories.UserTypeRepositoryInterface
}

type UserTypeServiceInterface interface {
}

func GetUserTypeService(userType repositories.UserTypeRepositoryInterface,
	log *log.Logger, cfg *config.Config) UserTypeServiceInterface {

	return &UserTypeService{
		log:          log,
		cfg:          cfg,
		userTypeRepo: userType,
	}
}

func (uts *UserTypeService) GetUserTypesList() ([]*models.UserType, error) {
	return uts.userTypeRepo.GetUserTypesList()
}
