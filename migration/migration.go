package migration

import (
	"public-surf/internal/domain/entity"
	"public-surf/pkg/database"
)

func Migrate() {
	db, err := database.InitDB()
	if err != nil {
		panic(err)
	}
	// user
	db.AutoMigrate(&entity.User{})

	// photo
	db.AutoMigrate(&entity.Photo{})
}
