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
	// Define the claims
	// claims := jwt.MapClaims{
	// 	"sub":   user.ID, // Subject identifier (replace with an actual identifier)
	// 	"name":  user.FirstName + " " + user.LastName,
	// 	"email": user.Email,
	// 	"exp":   time.Now().Add(time.Hour * 72).Unix(), // Token expiration time (1 hour in this example)
	// 	"iat":   time.Now().Unix(),                     // Issued at time
	// }

	// // Create a new token and specify the signing method
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// // Replace "your-secret-key" with your actual secret key
	// config := config.NewConfig()
	// secretKey := []byte(config.Jwt.SecretKey)

	// // Sign the token with the secret key
	// tokenString, err := token.SignedString(secretKey)
	// if err != nil {
	// 	return "", err
	// }
	jwttoken, err := jwttoken.CreateToken(user.Email)
	if err != nil {
		return "", err
	}
	tokenString := jwttoken.AccessToken

	return tokenString, nil
}
