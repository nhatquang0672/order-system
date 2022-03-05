package dao

import (
	"errors"
	"order-system/model"

	"gorm.io/gorm"
)

type OrderDAO struct {
	db *gorm.DB
}

func NewOrderDAO(db *gorm.DB) *OrderDAO {
	return &OrderDAO{
		db: db,
	}
}

func (od *OrderDAO) CreateCartItem(ct *model.CartItem) (err error) {
	if err := od.db.Create(&ct).Error; err != nil {
		return err
	}
	if err := od.db.Where(ct.ID).Preload("Product").First(&ct).Error; err != nil {
		return err
	}
	return nil
}

func (od *OrderDAO) UpdateCartItem(ct *model.CartItem) error {
	return od.db.Model(ct).Updates(ct).Error
}

func (od *OrderDAO) DeleteCartItem(ct *model.CartItem) error {
	return od.db.Delete(ct).Error
}

func (od *OrderDAO) CartItemFromId(id uint) (*model.CartItem, error) {
	var m model.CartItem
	if err := od.db.Preload("Product").First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (od *OrderDAO) ListCartItems(uid uint) ([]model.CartItem, error) {
	var cartitems []model.CartItem
	od.db.Where(&model.CartItem{UserID: uid}).Preload("Product").Find(&cartitems)
	return cartitems, nil
}

func (od *OrderDAO) PlaceOrder(uid uint, cartItems *[]model.CartItem) (*model.Order, error) {

	items := make([]model.OrderItem, 0)
	var totalPrice uint = 0
	for _, ci := range *cartItems {
		items = append(items, *model.OderItemFromCartItem(&ci))
		totalPrice += ci.Quantity * ci.Product.Price
	}
	o := &model.Order{UserID: uid, Status: "Paid", TotalPrice: totalPrice, Items: items}

	tx := od.db.Begin()
	if err := tx.Create(&o).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Delete(&cartItems).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Where(o.ID).Preload("Items.Product").First(&o).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return o, tx.Commit().Error
}

func (od *OrderDAO) GetOrderById(oid uint) (*model.Order, error) {
	var m model.Order
	err := od.db.Preload("Items.Product").First(&m, oid).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (od *OrderDAO) Orders(uid uint, offset, limit int) ([]model.Order, error) {
	var orders []model.Order
	od.db.Where(&model.Order{UserID: uid}).Preload("Items.Product").Offset(offset).Limit(limit).Order(" created_at desc").Find(&orders)
	return orders, nil
}
