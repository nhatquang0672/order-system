package db

import (
	"fmt"
	"order-system/model"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init() *gorm.DB {

	db, err := gorm.Open(sqlite.Open("./gorm.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("db err: ", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("db err: ", err)
	}

	sqlDB.SetMaxIdleConns(3)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}

func TestDBInit() *gorm.DB {
	test_db, err := gorm.Open(sqlite.Open("./../gorm_test.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("db err: ", err)
	}
	sqlDB, err := test_db.DB()
	if err != nil {
		fmt.Println("db err: ", err)
	}
	sqlDB.SetMaxIdleConns(3)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return test_db
}

func DropTestDB() error {
	if err := os.Remove("./../gorm_test.db"); err != nil {
		return err
	}
	return nil
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.User{},
		&model.Product{},
		&model.Order{},
		&model.OrderItem{},
		&model.CartItem{},
	)
}
