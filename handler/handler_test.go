package handler

import (
	"encoding/json"
	"log"
	"order-system/dao"
	"order-system/db"
	"order-system/model"
	"order-system/router"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var (
	d  *gorm.DB
	us *dao.UserDAO
	pd *dao.ProductDAO
	od *dao.OrderDAO
	h  *Handler
	e  *fiber.App
)

func tearDown() {
	if err := db.DropTestDB(); err != nil {
		log.Fatal(err)
	}
}

func setup() {
	d = db.TestDBInit()
	db.AutoMigrate(d)
	us = dao.NewUserDAO(d)
	pd = dao.NewProductDAO(d)
	od = dao.NewOrderDAO(d)
	h = NewHandler(us, pd, od)
	e = router.New()
	loadFakeData()
}

func responseMap(b []byte, key string) map[string]interface{} {
	var m map[string]interface{}
	json.Unmarshal(b, &m)
	return m[key].(map[string]interface{})
}

func authHeader(token string) string {
	return "Token " + token
}

func loadFakeData() error {
	u1 := model.User{
		Username: "user1",
		Email:    "user1@realworld.io",
		Role:     "user",
	}
	u1.Password, _ = u1.HashPassword("secret")
	if err := us.Create(&u1); err != nil {
		return err
	}

	u2 := model.User{
		Username: "user2",
		Email:    "user2@realworld.io",
		Role:     "vendor",
	}
	u2.Password, _ = u2.HashPassword("secret")
	if err := us.Create(&u2); err != nil {
		return err
	}

	p1 := model.Product{
		Description: "product 1 description",
		Price:       100,
		Quantity:    10,
	}
	pd.Create(&p1)

	p2 := model.Product{
		Description: "product 2 description",
		Price:       200,
		Quantity:    20,
	}
	pd.Create(&p2)

	ci1 := model.CartItem{
		ProductID: 2,
		UserID:    u1.ID,
		Quantity:  3,
	}
	od.CreateCartItem(&ci1)

	return nil
}
