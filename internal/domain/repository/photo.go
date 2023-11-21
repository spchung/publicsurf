package repository

import (
	"public-surf/internal/domain/entity"

	"gorm.io/gorm"
)

type PhotoRepository struct {
	db *gorm.DB
}

type IPhotoRepository interface {
	// FindAll() ([]entity.Photo, error)
	FindByID(id uint64) (entity.Photo, error)
	FindByUserID(userID string) ([]entity.Photo, error)
	FindLatestByUserID(userID uint64) (entity.PhotoViewModel, error)
	Save(photo *entity.Photo) (*entity.Photo, error)
	// Update(photo entity.Photo) (entity.Photo, error)
	// Delete(photo entity.Photo) (bool, error)
}

func NewPhotoRepository(db *gorm.DB) *PhotoRepository {
	return &PhotoRepository{db: db}
}

func NewPhoto() *entity.Photo {
	return &entity.Photo{}
}

func (r *PhotoRepository) FindByID(id uint64) (entity.Photo, error) {
	var photo entity.Photo
	err := r.db.First(&photo, id).Error
	return photo, err
}

func (r *PhotoRepository) FindByUserID(userID string) ([]entity.Photo, error) {
	var photos []entity.Photo
	err := r.db.Where("user_id = ?", userID).Find(&photos).Error
	return photos, err
}

func (r *PhotoRepository) FindLatestByUserID(userID uint64) (entity.PhotoViewModel, error) {
	var photo entity.PhotoViewModel
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").First(&photo).Error
	return photo, err
}

func (r *PhotoRepository) Save(photo *entity.Photo) (*entity.Photo, error) {
	err := r.db.Create(&photo).Error
	return photo, err
}
