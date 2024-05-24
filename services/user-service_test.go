package services

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"ypeskov/go_hillel_9/internal/config"
	"ypeskov/go_hillel_9/internal/log"
	"ypeskov/go_hillel_9/repository/models"
	"ypeskov/go_hillel_9/repository/repositories/mocks"
)

func TestCreateUser(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryInterface)
	mockCfg, _ := config.NewConfig()
	mockLog := log.New(mockCfg)

	service := GetUserService(mockRepo, mockLog, mockCfg)

	tests := []struct {
		name        string
		firstName   string
		lastName    string
		email       string
		password    string
		expectedErr error
		mockReturn  func()
	}{
		{
			name:        "TestCreateUser",
			firstName:   "Test",
			lastName:    "User",
			email:       "example@example.com",
			password:    "1",
			expectedErr: nil,
			mockReturn: func() {
				mockRepo.On("CreateUser", mock.MatchedBy(func(user *models.User) bool {
					return user.FirstName == "Test" &&
						user.LastName == "User" &&
						user.Email == "example@example.com" &&
						user.PasswordHash != ""
				})).Return(&models.User{
					Id:        1,
					FirstName: "Test",
					LastName:  "User",
					Email:     "example@example.com",
				}, nil)
			},
		},
		{
			name:        "TestCreateUser Failed Password Hash Empty String",
			firstName:   "Test",
			lastName:    "User",
			email:       "example@example.com",
			password:    "",
			expectedErr: errors.New("password too short to be a bcrypted password"),
			mockReturn: func() {
				mockRepo.On("CreateUser", mock.Anything).Return(nil, errors.New("password too short to be a bcrypted password"))
			},
		},
		{
			name:        "TestCreateUser Failed Password Hash Long String",
			firstName:   "Test",
			lastName:    "User",
			email:       "example@example.com",
			password:    string(make([]byte, 100)),
			expectedErr: errors.New("bcrypt: password length exceeds 72 bytes"),
			mockReturn: func() {
				mockRepo.On("CreateUser", mock.Anything).Return(nil, errors.New("bcrypt: password length exceeds 72 bytes"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockReturn()
			user, err := service.CreateUser(&models.User{
				FirstName:    tt.firstName,
				LastName:     tt.lastName,
				Email:        tt.email,
				PasswordHash: tt.password,
			})
			fmt.Printf("user: %v\n", user)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.firstName, user.FirstName)
				assert.Equal(t, tt.lastName, user.LastName)
				assert.Equal(t, tt.email, user.Email)
				assert.Greater(t, user.Id, 0)
			}
		})
	}
}

func TestGetUsersList(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryInterface)
	mockCfg, _ := config.NewConfig()
	mockLog := log.New(mockCfg)

	service := GetUserService(mockRepo, mockLog, mockCfg)

	tests := []struct {
		name       string
		mockReturn func()
	}{
		{
			name: "TestGetUsersList",
			mockReturn: func() {
				mockRepo.On("GetUsersList").Return([]*models.User{
					{
						Id:        1,
						FirstName: "Test",
						LastName:  "User",
						Email:     "example@example.com",
					},
				}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockReturn()
			users, err := service.GetUsersList()
			assert.NoError(t, err)
			assert.Len(t, users, 1)
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryInterface)
	mockCfg, _ := config.NewConfig()
	mockLog := log.New(mockCfg)

	service := GetUserService(mockRepo, mockLog, mockCfg)

	tests := []struct {
		name       string
		email      string
		mockReturn func()
	}{
		{
			name:  "TestGetUserByEmail",
			email: "example@example.com",
			mockReturn: func() {
				mockRepo.On("GetUserByEmail", "example@example.com").Return(&models.User{
					Id:        1,
					FirstName: "Test",
					LastName:  "User",
					Email:     "example@example.com",
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockReturn()
			user := service.GetUserByEmail(tt.email)
			assert.NotNil(t, user)
			assert.Equal(t, tt.email, user.Email)
		})
	}
}
