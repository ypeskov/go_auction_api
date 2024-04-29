package repositories

import (
	"ypeskov/go_hillel_9/internal/database"
	"ypeskov/go_hillel_9/internal/log"
	"ypeskov/go_hillel_9/repository/models"
)

type ItemRepository struct {
	log *log.Logger
	db  database.Database
}

func NewItemRepository(log *log.Logger, connection database.Database) *ItemRepository {
	return &ItemRepository{
		log: log,
		db:  connection,
	}
}

func (r *ItemRepository) GetItemsList() ([]*models.Item, error) {
	var items []*models.Item

	err := r.db.Select(&items, "SELECT * FROM items")
	if err != nil {
		r.log.Error("failed to get items from db", err)
		return nil, err
	}

	return items, nil
}

func (r *ItemRepository) CreateItem(srcItem *models.Item) (*models.Item, error) {
	insertQuery := "INSERT INTO items (user_id, title, initial_price, description) VALUES ($1, $2, $3, $4) RETURNING *"
	row, err := r.db.Query(insertQuery, srcItem.UserId, srcItem.Title, srcItem.InitialPrice, srcItem.Description)
	if err != nil {
		r.log.Error("failed to insert srcItem into db", err)
		return nil, err
	}

	var newItem models.Item
	if row.Next() {
		err = row.Scan(&newItem.Id, &newItem.UserId, &newItem.Title, &newItem.InitialPrice,
			&newItem.SoldPrice, &newItem.Description)
		if err != nil {
			r.log.Errorf("Failed to scan id: %v", err)
			return nil, err
		}
	}
	if newItem.SoldPrice == nil {
		var zero float64 = 0.0
		newItem.SoldPrice = &zero
	}

	return &newItem, nil
}
