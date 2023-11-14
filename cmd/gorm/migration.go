package main

import (
	"log"
	"public-surf/internal/domain/entity"
	"public-surf/pkg/config"
	"public-surf/pkg/database"
)

func init() {
	config.GetConfig()
}

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	// user
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.UserType{})
	// db.Model(&entity.User{}).AddForeignKey("user_type_id", "user_types(id)", "CASCADE", "CASCADE")

	// photo
	db.AutoMigrate(&entity.Photo{})
	db.AutoMigrate(&entity.PhotoFolder{})

	// cart
	db.AutoMigrate(&entity.Cart{})
	db.AutoMigrate(&entity.CartPhoto{})

	// order
	db.AutoMigrate(&entity.Order{})
	db.AutoMigrate(&entity.OrderPhoto{})
}
