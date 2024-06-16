package services

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"ypeskov/go_hillel_9/internal/config"
	"ypeskov/go_hillel_9/internal/log"
	"ypeskov/go_hillel_9/repository/models"
	"ypeskov/go_hillel_9/repository/repositories/mocks"
)

func TestGetItemsList(t *testing.T) {
	mockRepo := new(mocks.ItemRepositoryInterface)
	mockUserTypeRepo := new(mocks.UserTypeRepositoryInterface)
	mockCfg, _ := config.NewConfig()
	mockLog := log.New(mockCfg)

	service := GetItemService(mockRepo, mockUserTypeRepo, mockLog, mockCfg)

	userId := 1
	expectedItems := []*models.Item{
		{
			Id:           1,
			UserId:       1,
			Title:        "Test Item",
			InitialPrice: 100.0,
			SoldPrice:    nil,
			Description:  nil},
		{
			Id:           2,
			UserId:       1,
			Title:        "Test Item",
			InitialPrice: 100.0,
			SoldPrice:    nil,
			Description:  nil},
	}

	mockRepo.On("GetItemsList", userId).Return(expectedItems, nil)

	items, err := service.GetItemsList(userId)
	assert.NoError(t, err)
	assert.Equal(t, expectedItems, items)
	mockRepo.AssertExpectations(t)
}

func TestCreateItem(t *testing.T) {
	mockRepo := new(mocks.ItemRepositoryInterface)
	mockUserTypeRepo := new(mocks.UserTypeRepositoryInterface)
	mockCfg, _ := config.NewConfig()
	mockLog := log.New(mockCfg)

	service := GetItemService(mockRepo, mockUserTypeRepo, mockLog, mockCfg)

	srcItem := &models.Item{
		UserId:       1,
		Title:        "Test Item",
		InitialPrice: 100.0,
		SoldPrice:    nil,
		Description:  nil}
	expectedItem := &models.Item{
		Id:           1,
		UserId:       1,
		Title:        "Test Item",
		InitialPrice: 100.0,
		SoldPrice:    nil,
		Description:  nil}

	user := &models.User{
		Id:         1,
		FirstName:  "Test",
		LastName:   "User",
		Email:      "example@example.com",
		UserTypeId: 1,
	}

	mockRepo.On("CreateItem", srcItem, user).Return(expectedItem, nil)

	item, err := service.CreateItem(srcItem, user)
	fmt.Printf("%+v\n", item)
	assert.NoError(t, err)
	assert.Equal(t, expectedItem, item)
	mockRepo.AssertExpectations(t)
}

func TestGetItemById(t *testing.T) {
	mockRepo := new(mocks.ItemRepositoryInterface)
	mockUserTypeRepo := new(mocks.UserTypeRepositoryInterface)
	mockCfg, _ := config.NewConfig()
	mockLog := log.New(mockCfg)

	service := GetItemService(mockRepo, mockUserTypeRepo, mockLog, mockCfg)

	tests := []struct {
		name         string
		itemID       int
		userID       int
		expectedItem *models.Item
		expectedErr  error
		mockReturn   func()
	}{
		{
			name:   "Item found",
			itemID: 1,
			userID: 1,
			expectedItem: &models.Item{
				Id:           1,
				UserId:       1,
				Title:        "Test Item",
				InitialPrice: 100.0,
				SoldPrice:    nil,
				Description:  nil,
			},
			expectedErr: nil,
			mockReturn: func() {
				mockRepo.On("GetItemById", 1, 1).Return(&models.Item{
					Id:           1,
					UserId:       1,
					Title:        "Test Item",
					InitialPrice: 100.0,
					SoldPrice:    nil,
					Description:  nil,
				}, nil)
			},
		},
		{
			name:         "Item not found",
			itemID:       2,
			userID:       1,
			expectedItem: nil,
			expectedErr:  errors.New("item not found"),
			mockReturn: func() {
				mockRepo.On("GetItemById", 2, 1).Return(nil, errors.New("item not found"))
			},
		},
		{
			name:         "Database error",
			itemID:       1,
			userID:       2,
			expectedItem: nil,
			expectedErr:  errors.New("database error"),
			mockReturn: func() {
				mockRepo.On("GetItemById", 1, 2).Return(nil, errors.New("database error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockReturn()

			item, err := service.GetItemById(tt.itemID, tt.userID)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedItem, item)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateItem(t *testing.T) {
	mockRepo := new(mocks.ItemRepositoryInterface)
	mockUserTypeRepo := new(mocks.UserTypeRepositoryInterface)
	mockCfg, _ := config.NewConfig()
	mockLog := log.New(mockCfg)

	service := GetItemService(mockRepo, mockUserTypeRepo, mockLog, mockCfg)

	itemID := 1
	userID := 1
	srcItem := &models.Item{
		Id:           1,
		UserId:       1,
		Title:        "Test Item",
		InitialPrice: 100.0,
		SoldPrice:    nil,
		Description:  nil,
	}
	expectedItem := &models.Item{
		Id:           1,
		UserId:       1,
		Title:        "Test Item",
		InitialPrice: 100.0,
		SoldPrice:    nil,
		Description:  nil,
	}

	mockRepo.On("UpdateItem", itemID, srcItem, userID).Return(expectedItem, nil)

	item, err := service.UpdateItem(itemID, srcItem, userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedItem, item)
	mockRepo.AssertExpectations(t)
}

func TestDeleteItem(t *testing.T) {
	mockRepo := new(mocks.ItemRepositoryInterface)
	mockUserTypeRepo := new(mocks.UserTypeRepositoryInterface)
	mockCfg, _ := config.NewConfig()
	mockLog := log.New(mockCfg)

	service := GetItemService(mockRepo, mockUserTypeRepo, mockLog, mockCfg)

	itemID := 1
	userID := 1

	mockRepo.On("DeleteItem", itemID, userID).Return(nil)

	err := service.DeleteItem(itemID, userID)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
