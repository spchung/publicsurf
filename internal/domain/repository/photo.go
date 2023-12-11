package repository

import (
	"public-surf/internal/domain/entity"

	"gorm.io/gorm"
)

type PhotoRepository struct {
	db *gorm.DB
}

type IPhotoRepository interface {
	GetByID(id uint64) (entity.PhotoView, error)
	FindByUserID(userID uint64) ([]*entity.Photo, error)
	FindLatestByUserID(userID uint64) (entity.PhotoViewModel, error)
	Save(photo *entity.Photo) (*entity.Photo, error)
}

func NewPhotoRepository(db *gorm.DB) *PhotoRepository {
	return &PhotoRepository{db: db}
}

func NewPhoto() *entity.Photo {
	return &entity.Photo{}
}

func (r *PhotoRepository) GetByID(id uint64) (entity.PhotoView, error) {
	var photo entity.PhotoView
	err := r.db.First(&photo, id).Error
	return photo, err
}

func (r *PhotoRepository) FindByUserID(userID uint64) ([]*entity.Photo, error) {
	var photos []*entity.Photo
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
