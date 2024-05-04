package repositories

import (
	"fmt"
	"time"
	"ypeskov/go_hillel_9/internal/database"
	"ypeskov/go_hillel_9/internal/log"
	"ypeskov/go_hillel_9/repository/models"
)

type UserRepository struct {
	log *log.Logger
	db  database.Database
}

type UserRepositoryInterface interface {
	GetUsersList() ([]*models.User, error)
	CreateUser(srcUser *models.User) (*models.User, error)
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
	srcUser.LastLoginUtc = now

	insertQuery := `INSERT INTO users (first_name, last_name, email, password_hash, last_login_utc) 
                    VALUES (:first_name, :last_name, :email, :password_hash, :last_login_utc) RETURNING *`

	rows, err := r.db.NamedQuery(insertQuery, srcUser)
	if err != nil {
		r.log.Error("failed to insert srcUser into db", err)
		return nil, err
	}

	var newUser models.User
	if rows.Next() {
		err := rows.StructScan(&newUser)
		if err != nil {
			r.log.Errorf("Failed to scan user: %v", err)
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("failed to scan new user")
	}

	return &newUser, nil
}
