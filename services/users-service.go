package services

import (
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
	"ypeskov/go_hillel_9/internal/config"
	"ypeskov/go_hillel_9/internal/errors"
	"ypeskov/go_hillel_9/internal/log"
	"ypeskov/go_hillel_9/repository/models"
	"ypeskov/go_hillel_9/repository/repositories"
)

type UsersService struct {
	log      *log.Logger
	cfg      *config.Config
	userRepo repositories.UserRepositoryInterface
}

type Claims struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type UsersServiceInterface interface {
	CreateUser(srcUser *models.User) (*models.User, error)
	GetUsersList() ([]*models.User, error)
	GetJWT(email string, password string) (string, error)
}

func GetUserService(userRepo repositories.UserRepositoryInterface,
	log *log.Logger, cfg *config.Config) UsersServiceInterface {

	return &UsersService{
		log:      log,
		cfg:      cfg,
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

func (us *UsersService) GetUsersList() ([]*models.User, error) {
	return us.userRepo.GetUsersList()
}

func (us *UsersService) GetJWT(username string, password string) (string, error) {
	user := us.userRepo.GetUserByEmail(username)
	if user == nil {
		return "", errors.NotFoundErr
	}

	if !checkPasswordHash(password, user.PasswordHash) {
		return "", errors.UnauthorizedErr
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		Id:    user.Id,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(us.cfg.SECRET_KEY))
	if err != nil {
		us.log.Errorln("failed to sign token: ", err)
		return "", errors.InternalServerErr
	}

	return tokenString, nil
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
