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
	FindByID(id uint64) (entity.User, error)
	FindByEmail(email string) (entity.User, error)
	// Store(user entity.User) (entity.User, error)
	// Update(user entity.User) (entity.User, error)
	// Delete(user entity.User) (bool, error)
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(id uint64) (entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	return user, err
}

func (r *UserRepository) FindByEmail(email string) (entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return user, err
}
