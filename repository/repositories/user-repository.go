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
	GetUserByEmail(email string) *models.User
	AddOrUpdateRefreshToken(userId int, token string) error
	GetUserByRefreshToken(token string) *models.User
}

func GetUserRepository(log *log.Logger, connection database.Database) UserRepositoryInterface {
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

func (r *UserRepository) GetUserByEmail(email string) *models.User {
	var user models.User

	err := r.db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		r.log.Error("failed to get user from db", err)

		return nil
	}

	return &user
}

func (r *UserRepository) AddOrUpdateRefreshToken(userId int, token string) error {
	query := `
        INSERT INTO refresh_tokens (user_id, token, created_at)
        VALUES ($1, $2, now())
        ON CONFLICT (user_id) DO UPDATE
        SET token = EXCLUDED.token, created_at = now()
    `
	_, err := r.db.Exec(query, userId, token)
	if err != nil {
		r.log.Error("failed to insert or update refresh token in db", err)

		return err
	}

	return nil
}

func (r *UserRepository) GetUserByRefreshToken(token string) *models.User {
	var user models.User

	err := r.db.Get(&user,
		`SELECT u.* FROM users u JOIN refresh_tokens rt ON u.id = rt.user_id WHERE rt.token = $1`, token)
	if err != nil {
		r.log.Errorln("failed to get user from db with refresh token", err)

		return nil
	}

	return &user
}
