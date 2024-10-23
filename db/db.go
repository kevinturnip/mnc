package db

import (
	"mnc/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := "host=localhost user=user1 password=password dbname=mnc port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	DB.AutoMigrate(&model.User{}, &model.Transaction{})
}
