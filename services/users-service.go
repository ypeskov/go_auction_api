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
	GetJWT(email string, password string, minutes time.Duration, refreshToken bool) (string, error)
	GetRefreshToken(email string, password string, update bool) (string, error)
	GetUserByEmail(email string) *models.User
	GetUserByRefreshToken(token string) (*models.User, error)
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

func (us *UsersService) GetJWT(email string, password string, minutes time.Duration, refreshToken bool) (string, error) {
	user := us.userRepo.GetUserByEmail(email)
	if user == nil {
		return "", errors.NotFoundErr
	}

	// If we are not refreshing the token, we need to check the password
	if refreshToken == false {
		if !checkPasswordHash(password, user.PasswordHash) {
			return "", errors.UnauthorizedErr
		}
	}

	expirationTime := time.Now().Add(minutes * time.Minute)

	claims := &Claims{
		Id:    user.Id,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(us.cfg.SecretKey))
	if err != nil {
		us.log.Errorln("failed to sign token: ", err)
		return "", errors.InternalServerErr
	}

	return tokenString, nil
}

func (us *UsersService) GetRefreshToken(email string, password string, update bool) (string, error) {
	user := us.userRepo.GetUserByEmail(email)
	if user == nil {
		return "", errors.NotFoundErr
	}

	minutes := time.Duration(us.cfg.RefreshTokenLifetimeMinutes)
	token, err := us.GetJWT(email, password, minutes, update)
	if err != nil {
		us.log.Errorln("failed to get JWT", err)
		return "", err
	}

	err = us.userRepo.AddOrUpdateRefreshToken(user.Id, token)
	if err != nil {
		us.log.Errorln("failed to add refresh token", err)
		return "", errors.InternalServerErr
	}

	return token, nil
}

func (us *UsersService) GetUserByEmail(email string) *models.User {
	return us.userRepo.GetUserByEmail(email)
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

func (us *UsersService) GetUserByRefreshToken(token string) (*models.User, error) {
	user := us.userRepo.GetUserByRefreshToken(token)
	if user == nil {
		us.log.Errorln("failed to get user by refresh token")
		return nil, errors.UnauthorizedErr
	}

	return user, nil
}
