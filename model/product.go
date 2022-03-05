package model

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Description string
	Quantity    uint
	Price       uint
}

type Order struct {
	gorm.Model
	TotalPrice uint
	Status     string
	User       User
	UserID     uint
	Items      []OrderItem
}

type OrderItem struct {
	gorm.Model
	Order     Order
	OrderID   uint
	Product   Product
	ProductID uint
	Quantity  uint
	Price     uint
}

type CartItem struct {
	gorm.Model
	User      User
	UserID    uint
	Product   Product
	ProductID uint
	Quantity  uint
}

func OderItemFromCartItem(ci *CartItem) *OrderItem {
	oi := &OrderItem{}
	oi.ProductID = ci.ProductID
	oi.Quantity = ci.Quantity
	oi.Product = ci.Product
	oi.Price = ci.Quantity * ci.Product.Price
	return oi
}
