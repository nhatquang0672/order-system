package main

import (
	"fmt"
	"order-system/dao"
	"order-system/db"
	"order-system/handler"
	"order-system/model"
	"order-system/router"

	_ "order-system/docs"

	swagger "github.com/arsmn/fiber-swagger/v2"
)

// @description Order System API
// @title Order System API

// @BasePath /api

// @schemes http https
// @produce application/json
// @consumes application/json

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	r := router.New()
	r.Get("/swagger/*", swagger.HandlerDefault)

	d := db.Init()
	db.AutoMigrate(d)

	us := dao.NewUserDAO(d)
	pd := dao.NewProductDAO(d)
	od := dao.NewOrderDAO(d)

	////////////////////////////////////////
	loadFakeData(us, pd)
	////////////////////////////////////////

	h := handler.NewHandler(us, pd, od)
	h.Register(r)
	err := r.Listen(":8585")
	if err != nil {
		fmt.Printf("%v", err)
	}
}

func loadFakeData(us *dao.UserDAO, pd *dao.ProductDAO) error {
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

	return nil
}
