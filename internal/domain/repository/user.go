package repository

import (
	"public-surf/internal/domain/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type IUserRepository interface {
	// FindAll() ([]entity.User, error)
	GetUser(id int) (entity.User, error)
	GetUserByEmail(email string) (entity.User, error)
	// Store(user entity.User) (entity.User, error)
	// Update(user entity.User) (entity.User, error)
	// Delete(user entity.User) (bool, error)
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUser(id int) (entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	return user, err
}

func (r *UserRepository) GetUserByEmail(email string) (entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return user, err
}
