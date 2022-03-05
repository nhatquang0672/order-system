package dao

import (
	"errors"
	"order-system/model"

	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (us *UserDAO) GetByID(id uint) (*model.User, error) {
	var m model.User
	if err := us.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (us *UserDAO) GetByEmail(e string) (*model.User, error) {
	var m model.User
	if err := us.db.Where(&model.User{Email: e}).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (us *UserDAO) Create(u *model.User) (err error) {
	return us.db.Create(u).Error
}

func (us *UserDAO) Update(u *model.User) error {
	return us.db.Model(u).Updates(u).Error
}
