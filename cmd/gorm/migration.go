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

	// photo
	db.AutoMigrate(&entity.Photo{})
	db.AutoMigrate(&entity.PhotoFolder{})

	// cart
	db.AutoMigrate(&entity.Cart{})
	db.AutoMigrate(&entity.CartPhoto{})

	// order
	db.AutoMigrate(&entity.Order{})
	db.AutoMigrate(&entity.OrderPhoto{})

	// photo view
	db.Exec(`
		CREATE OR REPLACE VIEW v_photos AS SELECT 
			sub.*,
			pt."name" as "user_type"
		from photo_types pt JOIN (
			SELECT 
			p.*,
			u.email as "user_email"
		FROM photos p JOIN users u ON p.user_id = u.id) as sub on sub.photo_type_id = pt.id
	`)
}
