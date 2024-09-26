package repositories

import (
	"errors"
	"first_gin_app/models"

	"gorm.io/gorm"
)

// itemのポインタを全取得するように宣言
type IItemRepository interface {
	FindAll() (*[]models.Item, error)
	FindById(itemId uint, userId uint) (*models.Item, error)
	Create(newItem models.Item) (*models.Item, error)
	Update(updateItem models.Item) (*models.Item, error)
	Delete(itemId uint, userId uint) error
}

// // Itemの配列を保持する属性を定義
// type ItemMemoryRepository struct {
// 	items []models.Item
// }

// // ItemMemoryRepositoryの参照を返す
// func NewItemMemoryRepository(items []models.Item) IItemRepository {
// 	return &ItemMemoryRepository{items: items}
// }

// // 自信(ItemMemoryRepository構造体のオブジェクト)がitems属性の値を全て返す
// func (r *ItemMemoryRepository) FindAll() (*[]models.Item, error) {
// 	return &r.items, nil
// }

// func (r *ItemMemoryRepository) FindById(itemId uint) (*models.Item, error) {
// 	for _, v := range r.items {
// 		if v.ID == itemId {
// 			return &v, nil
// 		}
// 	}
// 	return nil, errors.New("item not found")
// }

// func (r *ItemMemoryRepository) Create(newItem models.Item) (*models.Item, error) {
// 	newItem.ID = uint(len(r.items) + 1)
// 	r.items = append(r.items, newItem)
// 	return &newItem, nil
// }

// func (r *ItemMemoryRepository) Update(updateItem models.Item) (*models.Item, error) {
// 	for i, v := range r.items {
// 		if v.ID == updateItem.ID {
// 			r.items[i] = updateItem
// 			return &r.items[i], nil
// 		}
// 	}
// 	return nil, errors.New("unexpected error")
// }

// func (r *ItemMemoryRepository) Delete(itemId uint) error {
// 	for i, v := range r.items {
// 		if v.ID == itemId {
// 			r.items = append(r.items[:i], r.items[i+1:]...)
// 			return nil
// 		}
// 	}
// 	return errors.New("item not found")
// }

type ItemRepository struct {
	db *gorm.DB
}

// Create implements IItemRepository.
func (r *ItemRepository) Create(newItem models.Item) (*models.Item, error) {
	result := r.db.Create(&newItem)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newItem, nil
}

// Delete implements IItemRepository.
func (r *ItemRepository) Delete(itemId uint, userId uint) error {
	deleteItem, err := r.FindById(itemId, userId)
	if err != nil {
		return err
	}
	// デフォルトは論理削除；データ削除を行いたい場合はr.db.Unscoped().Delete(&deleteItem)
	result := r.db.Delete(&deleteItem)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindAll implements IItemRepository.
func (r *ItemRepository) FindAll() (*[]models.Item, error) {
	var items []models.Item
	result := r.db.Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return &items, nil
}

// FindById implements IItemRepository.
func (r *ItemRepository) FindById(itemId uint, userId uint) (*models.Item, error) {
	var item models.Item
	result := r.db.First(&item, "id = ? AND user_id = ?", itemId, userId)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("item not found")
		}
		return nil, result.Error
	}
	return &item, nil
}

// Update implements IItemRepository.
func (r *ItemRepository) Update(updateItem models.Item) (*models.Item, error) {
	result := r.db.Save(&updateItem)
	if result.Error != nil {
		return nil, result.Error
	}
	return &updateItem, nil
}

func NewItemRepository(db *gorm.DB) IItemRepository {
	return &ItemRepository{db: db}
}
