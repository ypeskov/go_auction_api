package repositories

import (
	"time"
	"ypeskov/go_hillel_9/internal/database"
	"ypeskov/go_hillel_9/internal/log"
	"ypeskov/go_hillel_9/repository/models"
)

type UserRepository struct {
	log *log.Logger
	db  database.Database
}

func GetUserRepository(log *log.Logger, connection database.Database) *UserRepository {
	return &UserRepository{
		log: log,
		db:  connection,
	}
}

func (r *UserRepository) GetUsersList() ([]*models.User, error) {
	var users []*models.User

	err := r.db.Select(&users, "SELECT * FROM users")
	if err != nil {
		r.log.Error("failed to get users from db", err)
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) CreateUser(srcUser *models.User) (*models.User, error) {
	now := time.Now().UTC()

	insertQuery := "INSERT INTO users (first_name, last_name, email, last_login_utc) VALUES ($1, $2, $3, $4) RETURNING *"
	row, err := r.db.Query(insertQuery, srcUser.FirsName, srcUser.LastName, srcUser.Email, now)
	if err != nil {
		r.log.Error("failed to insert srcUser into db", err)
		return nil, err
	}

	var newUser models.User
	if row.Next() {
		err = row.Scan(&newUser.Id, &newUser.FirsName, &newUser.LastName, &newUser.Email, &newUser.LastLoginUtc)
		if err != nil {
			r.log.Errorf("Failed to scan id: %v", err)
			return nil, err
		}
	} else {
		r.log.Error("no rows returned")
		return nil, err
	}

	return &newUser, nil
}
