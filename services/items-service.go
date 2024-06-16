package services

import (
	"mime/multipart"
	"ypeskov/go_hillel_9/internal/config"
	"ypeskov/go_hillel_9/internal/log"
	"ypeskov/go_hillel_9/repository/models"
	"ypeskov/go_hillel_9/repository/repositories"
)

type ItemService struct {
	log          *log.Logger
	cfg          *config.Config
	itemRepo     repositories.ItemRepositoryInterface
	userTypeRepo repositories.UserTypeRepositoryInterface
}

type ItemsServiceInterface interface {
	GetItemsList(userId int) ([]*models.Item, error)
	CreateItem(srcItem *models.Item, user *models.User) (*models.Item, error)
	GetItemById(id int, userId int) (*models.Item, error)
	UpdateItem(id int, srcItem *models.Item, userId int) (*models.Item, error)
	DeleteItem(id int, userid int) error
	GetAllItems() ([]*models.Item, error)
	CreateItemComment(comment *models.ItemComment) (*models.ItemComment, error)
}

func GetItemService(itemRepo repositories.ItemRepositoryInterface,
	userTypeRepo repositories.UserTypeRepositoryInterface,
	log *log.Logger, cfg *config.Config) ItemsServiceInterface {

	return &ItemService{
		log:          log,
		cfg:          cfg,
		itemRepo:     itemRepo,
		userTypeRepo: userTypeRepo,
	}
}

func (is *ItemService) GetItemsList(userId int) ([]*models.Item, error) {
	return is.itemRepo.GetItemsList(userId)
}

func (is *ItemService) GetAllItems() ([]*models.Item, error) {
	return is.itemRepo.GetAllItems()
}

func (is *ItemService) CreateItem(srcItem *models.Item, user *models.User) (*models.Item, error) {
	userTypes, err := is.userTypeRepo.GetUserTypesList()
	if err != nil {
		return nil, err
	}

	if !canUserAddItem(user, userTypes) {
		is.log.Errorf("User type is not SELLER: %+v\n", srcItem)
		return nil, IncorrectUserRoleErr

	}

	return is.itemRepo.CreateItem(srcItem)
}

func (is *ItemService) GetItemById(id int, userId int) (*models.Item, error) {
	return is.itemRepo.GetItemById(id, userId)
}

func (is *ItemService) UpdateItem(id int, srcItem *models.Item, userId int) (*models.Item, error) {
	return is.itemRepo.UpdateItem(id, srcItem, userId)
}

func (is *ItemService) DeleteItem(id int, userId int) error {
	return is.itemRepo.DeleteItem(id, userId)
}

func canUserAddItem(user *models.User, userTypes []*models.UserType) bool {
	var sellerTypeId int32
	for _, userType := range userTypes {
		if userType.TypeCode == "SELLER" {
			sellerTypeId = int32(userType.Id)
		}
	}

	return user.UserTypeId == sellerTypeId
}

func (is *ItemService) CreateItemComment(comment *models.ItemComment) (*models.ItemComment, error) {
	return is.itemRepo.CreateItemComment(comment)
}

func (is *ItemService) AttachFileToItem(itemId int, file *multipart.FileHeader) error {
	return nil
}
