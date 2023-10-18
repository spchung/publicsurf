package service

import (
	"public-surf/internal/domain/entity"
	"public-surf/internal/domain/repository"
	"strconv"
)

type UserService struct {
	userRepo  repository.IUserRepository
	photoRepo repository.IPhotoRepository
}

type IUserService interface {
	GetUserPhotoCount(id uint64) (int64, error)
	GetUserLatestPhoto(id uint64) (entity.PhotoViewModel, error)
	GetUser(id uint64) (entity.UserViewModel, error)
}

func NewUserService(userRepo repository.IUserRepository, photoRepo repository.IPhotoRepository) *UserService {
	return &UserService{userRepo: userRepo, photoRepo: photoRepo}
}

func (s *UserService) GetUserPhotoCount(userID uint64) (int64, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return 0, err
	}
	photos, err := s.photoRepo.FindByUserID(strconv.Itoa(int(user.ID)))
	if err != nil {
		return 0, err
	}
	return int64(len(photos)), nil
}

func (s *UserService) GetUser(id uint64) (entity.UserViewModel, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return entity.UserViewModel{}, err
	}
	return entity.UserViewModel{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}

func (s *UserService) GetUserLatestPhoto(userID uint64) (entity.PhotoViewModel, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return entity.PhotoViewModel{}, err
	}
	photo, err := s.photoRepo.FindLatestByUserID(user.ID)
	if err != nil {
		return entity.PhotoViewModel{}, err
	}
	return photo, nil
}
