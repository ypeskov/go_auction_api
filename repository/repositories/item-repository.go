package repositories

import (
	"fmt"
	"time"
	"ypeskov/go_hillel_9/internal/database"
	"ypeskov/go_hillel_9/internal/errors"
	"ypeskov/go_hillel_9/internal/log"
	"ypeskov/go_hillel_9/repository/models"
)

type ItemRepository struct {
	log *log.Logger
	db  database.Database
}

type ItemRepositoryInterface interface {
	GetItemsList(userId int) ([]*models.Item, error)
	CreateItem(srcItem *models.Item) (*models.Item, error)
	GetItemById(id int, userId int) (*models.Item, error)
	UpdateItem(id int, srcItem *models.Item, userId int) (*models.Item, error)
	DeleteItem(id int, userId int) error
	GetAllItems() ([]*models.Item, error)
	CreateItemComment(comment *models.ItemComment) (*models.ItemComment, error)
}

func GetItemRepository(log *log.Logger, connection database.Database) ItemRepositoryInterface {
	return &ItemRepository{
		log: log,
		db:  connection,
	}
}

func (r *ItemRepository) GetItemsList(userId int) ([]*models.Item, error) {
	var items []*models.Item

	err := r.db.Select(&items, "SELECT * FROM items WHERE user_id = $1", userId)
	if err != nil {
		r.log.Error("failed to get items from db", err)

		return nil, err
	}

	return items, nil
}

func (r *ItemRepository) CreateItem(srcItem *models.Item) (*models.Item, error) {
	insertQuery := `INSERT INTO items (user_id, title, initial_price, description) 
					VALUES ($1, $2, $3, $4) RETURNING *`
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

func (r *ItemRepository) GetItemById(id int, userId int) (*models.Item, error) {
	var item models.Item

	err := r.db.Get(&item, "SELECT * FROM items WHERE id = $1 AND user_id = $2", id, userId)
	if err != nil {
		r.log.Errorln("failed to get item by id", err)

		return nil, err
	}

	return &item, nil
}

func (r *ItemRepository) UpdateItem(id int, srcItem *models.Item, userId int) (*models.Item, error) {
	updateQuery :=
		"UPDATE items SET user_id = $1, title = $2, initial_price = $3, description = $4 " +
			"WHERE id = $5 AND user_id = $6 RETURNING *"
	row, err := r.db.Query(updateQuery, userId, srcItem.Title, srcItem.InitialPrice,
		srcItem.Description, id, userId)
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

func (r *ItemRepository) DeleteItem(id int, userId int) error {
	result, err := r.db.Exec("DELETE FROM items WHERE id = $1 AND user_id = $2", id, userId)
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

func (r *ItemRepository) GetAllItems() ([]*models.Item, error) {
	var items []*models.Item

	err := r.db.Select(&items, "SELECT * FROM items")
	if err != nil {
		r.log.Error("failed to get items from db", err)

		return nil, err
	}

	return items, nil
}

func (r *ItemRepository) CreateItemComment(comment *models.ItemComment) (*models.ItemComment, error) {
	now := time.Now()
	comment.CreatedAt = now

	insertQuery := `INSERT INTO item_comments (user_id, item_id, comment, created_at) 
					VALUES (:user_id, :item_id, :comment, :created_at) RETURNING *`

	rows, err := r.db.NamedQuery(insertQuery, comment)
	if err != nil {
		r.log.Error("failed to insert comment into db", err)
		r.log.Errorf("comment: %+v\n", comment)

		return nil, err
	}

	var newComment models.ItemComment
	if rows.Next() {
		err := rows.StructScan(&newComment)
		if err != nil {
			r.log.Errorf("Failed to scan comment: %v", err)

			return nil, err
		}
	} else {
		return nil, fmt.Errorf("failed to add a new comment")
	}

	return &newComment, nil
}
