package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"ypeskov/go_hillel_9/internal/config"
	"ypeskov/go_hillel_9/internal/log"
	"ypeskov/go_hillel_9/internal/utils"
	"ypeskov/go_hillel_9/repository/models"
	"ypeskov/go_hillel_9/repository/repositories"
)

type ItemService struct {
	log          *log.Logger
	cfg          *config.Config
	itemRepo     repositories.ItemRepositoryInterface
	userTypeRepo repositories.UserTypeRepositoryInterface
}

const uploadPath = "./uploads"

type Bid struct {
	ItemId int
	Amount float64
}

type ItemsServiceInterface interface {
	GetItemsList(userId int) ([]*models.Item, error)
	CreateItem(srcItem *models.Item, user *models.User) (*models.Item, error)
	GetItemById(id int, userId int) (*models.Item, error)
	UpdateItem(id int, srcItem *models.Item, userId int) (*models.Item, error)
	DeleteItem(id int, userid int) error
	GetAllItems() ([]*models.Item, error)
	CreateItemComment(comment *models.ItemComment) (*models.ItemComment, error)
	AttachFileToItem(itemId int, file *multipart.FileHeader) (*string, error)
	CreateBid(bidChannel chan<- Bid, itemId int, amount float64) error
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

func (is *ItemService) AttachFileToItem(itemId int, file *multipart.FileHeader) (*string, error) {
	src, err := file.Open()
	if err != nil {
		is.log.Error("failed to open file", err)

		return nil, err
	}
	defer src.Close()

	// Destination
	err = utils.EnsureDir(uploadPath)
	if err != nil {
		is.log.Errorln("failed to ensure dir for file", err)

		return nil, err
	}
	fileName := fmt.Sprintf("%d_%s", itemId, file.Filename)
	fullFileName := fmt.Sprintf("%s/%s", uploadPath, fileName)
	dst, err := os.Create(fullFileName)
	if err != nil {
		is.log.Errorln("failed to create file", err)

		return nil, err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		is.log.Errorln("failed to copy file", err)

		return nil, err
	}

	err = is.itemRepo.AttachFileToItem(itemId, fileName)
	if err != nil {
		is.log.Errorln("failed to attach file to item", err)

		return nil, err
	}

	return &fullFileName, nil
}

func (is *ItemService) CreateBid(bidChannel chan<- Bid, itemId int, amount float64) error {
	// TODO: add storing bid to DB

	bidChannel <- Bid{
		ItemId: itemId,
		Amount: amount,
	}

	return nil
}
