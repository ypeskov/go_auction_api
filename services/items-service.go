package services

import (
	"ypeskov/go_hillel_9/internal/config"
	"ypeskov/go_hillel_9/internal/log"
	"ypeskov/go_hillel_9/repository/models"
	"ypeskov/go_hillel_9/repository/repositories"
)

type ItemService struct {
	log      *log.Logger
	cfg      *config.Config
	itemRepo repositories.ItemRepositoryInterface
}

type ItemsServiceInterface interface {
	GetItemsList() ([]*models.Item, error)
	CreateItem(srcItem *models.Item) (*models.Item, error)
	GetItemById(id int) (*models.Item, error)
	UpdateItem(id int, srcItem *models.Item) (*models.Item, error)
	DeleteItem(id int) error
}

func GetItemService(itemRepo repositories.ItemRepositoryInterface,
	log *log.Logger, cfg *config.Config) ItemsServiceInterface {

	return &ItemService{
		log:      log,
		cfg:      cfg,
		itemRepo: itemRepo,
	}
}

func (is *ItemService) GetItemsList() ([]*models.Item, error) {
	return is.itemRepo.GetItemsList()
}

func (is *ItemService) CreateItem(srcItem *models.Item) (*models.Item, error) {
	return is.itemRepo.CreateItem(srcItem)
}

func (is *ItemService) GetItemById(id int) (*models.Item, error) {
	return is.itemRepo.GetItemById(id)
}

func (is *ItemService) UpdateItem(id int, srcItem *models.Item) (*models.Item, error) {
	return is.itemRepo.UpdateItem(id, srcItem)
}

func (is *ItemService) DeleteItem(id int) error {
	return is.itemRepo.DeleteItem(id)
}
