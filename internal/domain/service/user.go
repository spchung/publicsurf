package service

import (
	"public-surf/internal/domain/entity"
	"public-surf/internal/domain/repository"
)

type UserService struct {
	userRepo  repository.IUserRepository
	photoRepo repository.IPhotoRepository
}

type IUserService interface {
	GetUserPhotoCount(id int) (int, error)
	GetUserLatestPhoto(id int) (entity.PhotoViewModel, error)
	GetUser(id int) (entity.UserViewModel, error)
	GetUserByEmail(email string) (entity.User, error)
}

func NewUserService(userRepo repository.IUserRepository, photoRepo repository.IPhotoRepository) *UserService {
	return &UserService{userRepo: userRepo, photoRepo: photoRepo}
}

func (s *UserService) GetUserPhotoCount(userID int) (int, error) {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return 0, err
	}
	photos, err := s.photoRepo.FindByUserID(user.ID)
	if err != nil {
		return 0, err
	}
	return len(photos), nil
}

func (s *UserService) GetUser(id int) (entity.UserViewModel, error) {
	user, err := s.userRepo.GetUser(id)
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

func (s *UserService) GetUserLatestPhoto(userID int) (entity.PhotoViewModel, error) {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return entity.PhotoViewModel{}, err
	}
	photo, err := s.photoRepo.FindLatestByUserID(user.ID)
	if err != nil {
		return entity.PhotoViewModel{}, err
	}
	return photo, nil
}

func (s *UserService) GetUserByEmail(email string) (entity.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}
