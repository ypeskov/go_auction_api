package services

import (
	"golang.org/x/crypto/bcrypt"
	"ypeskov/go_hillel_9/internal/log"
	"ypeskov/go_hillel_9/repository/models"
	"ypeskov/go_hillel_9/repository/repositories"
)

type UsersService struct {
	log      *log.Logger
	userRepo repositories.UserRepositoryInterface
}

type UsersServiceInterface interface {
	CreateUser(srcUser *models.User) (*models.User, error)
	GetJWT(username string, password string) (string, error)
}

func GetUserService(userRepo repositories.UserRepositoryInterface, log *log.Logger) UsersServiceInterface {
	return &UsersService{
		log:      log,
		userRepo: userRepo,
	}
}

func (us *UsersService) CreateUser(srcUser *models.User) (*models.User, error) {
	hash, err := hashPassword(srcUser.PasswordHash)
	if err != nil {
		us.log.Error("failed to hash password", err)
		return nil, err
	}
	srcUser.PasswordHash = hash

	return us.userRepo.CreateUser(srcUser)
}

func (us *UsersService) GetJWT(username string, password string) (string, error) {
	// This is a placeholder for a real JWT generation
	return "jwt", nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
