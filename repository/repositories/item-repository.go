package repositories

import (
	"ypeskov/go_hillel_9/internal/database"
	"ypeskov/go_hillel_9/internal/errors"
	"ypeskov/go_hillel_9/internal/log"
	"ypeskov/go_hillel_9/repository/models"
)

type ItemRepository struct {
	log *log.Logger
	db  database.Database
}

func GetItemRepository(log *log.Logger, connection database.Database) *ItemRepository {
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
	} else {
		r.log.Error("failed to scan new item")
		return nil, err

	}
	if newItem.SoldPrice == nil {
		var zero float64 = 0.0
		newItem.SoldPrice = &zero
	}

	return &newItem, nil
}

func (r *ItemRepository) GetItemById(id int) (*models.Item, error) {
	var item models.Item

	err := r.db.Get(&item, "SELECT * FROM items WHERE id = $1", id)
	if err != nil {
		r.log.Error("failed to get item by id", err)
		return nil, err
	}

	return &item, nil
}

func (r *ItemRepository) UpdateItem(id int, srcItem *models.Item) (*models.Item, error) {
	updateQuery :=
		"UPDATE items SET user_id = $1, title = $2, initial_price = $3, description = $4 WHERE id = $5 RETURNING *"
	row, err := r.db.Query(updateQuery, srcItem.UserId, srcItem.Title, srcItem.InitialPrice, srcItem.Description, id)
	if err != nil {
		r.log.Error("failed to update item in db", err)
		return nil, err
	}

	var updatedItem models.Item
	if row.Next() {
		err = row.Scan(&updatedItem.Id, &updatedItem.UserId, &updatedItem.Title, &updatedItem.InitialPrice,
			&updatedItem.SoldPrice, &updatedItem.Description)
		if err != nil {
			r.log.Errorf("Failed to scan id: %v", err)
			return nil, err
		}
	}
	if updatedItem.SoldPrice == nil {
		var zero float64 = 0.0
		updatedItem.SoldPrice = &zero
	}

	return &updatedItem, nil
}

func (r *ItemRepository) DeleteItem(id int) error {
	result, err := r.db.Exec("DELETE FROM items WHERE id = $1", id)
	if err != nil {
		r.log.Error("failed to delete item from db", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.log.Error("error checking rows affected", err)
		return err
	}

	if rowsAffected == 0 {
		r.log.Errorln("no item was deleted")
		return errors.NotFoundErr
	} else {
		r.log.Infof("Item with id %d deleted", id)
	}

	return nil
}
