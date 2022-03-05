package dao

import (
	"errors"
	"order-system/model"

	"gorm.io/gorm"
)

type ProductDAO struct {
	db *gorm.DB
}

func NewProductDAO(db *gorm.DB) *ProductDAO {
	return &ProductDAO{
		db: db,
	}
}

func (pd *ProductDAO) GetById(id uint) (*model.Product, error) {
	var m model.Product
	if err := pd.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (pd *ProductDAO) Create(u *model.Product) (err error) {
	return pd.db.Create(u).Error
}

func (pd *ProductDAO) Update(u *model.Product) error {
	return pd.db.Model(u).Updates(u).Error
}

func (pd *ProductDAO) List(offset, limit int) ([]model.Product, int64, error) {
	var (
		products []model.Product
		count    int64
	)
	pd.db.Model(&products).Count(&count)
	pd.db.Limit(limit).Offset(offset).Order(" created_at desc").Find(&products)
	return products, count, nil
}
