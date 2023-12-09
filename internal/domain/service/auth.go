package service

import (
	"public-surf/internal/domain/entity"
	"public-surf/internal/domain/repository"
	"public-surf/pkg/jwttoken"
)

type AuthService struct {
	userRepo repository.IUserRepository
}

type IAuthService interface {
	GenerateJwt(user *entity.User) (string, error)
}

func NewAuthService(userRepo repository.IUserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (svc *AuthService) GenerateJwt(user *entity.User) (string, error) {
	jwttoken, err := jwttoken.CreateToken(user.Email)
	if err != nil {
		return "", err
	}
	tokenString := jwttoken.AccessToken

	return tokenString, nil
}
