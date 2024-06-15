package repositories

import (
	"ypeskov/go_hillel_9/internal/database"
	"ypeskov/go_hillel_9/internal/log"
	"ypeskov/go_hillel_9/repository/models"
)

type UserTypeRepository struct {
	log *log.Logger
	db  database.Database
}

type UserTypeRepositoryInterface interface {
	GetUserTypesList() ([]*models.UserType, error)
}

func GetUserTypeRepository(log *log.Logger, connection database.Database) UserTypeRepositoryInterface {
	return &UserTypeRepository{
		log: log,
		db:  connection,
	}
}

func (r *UserTypeRepository) GetUserTypesList() ([]*models.UserType, error) {
	var userTypes []*models.UserType

	err := r.db.Select(&userTypes, "SELECT * FROM user_types")
	if err != nil {
		r.log.Errorln("failed to get user_types from db", err)

		return nil, err
	}

	return userTypes, nil
}
