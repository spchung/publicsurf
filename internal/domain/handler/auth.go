package handler

import (
	"net/http"
	"public-surf/internal/domain/service"
	"public-surf/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.IAuthService
	userService service.IUserService
}

func NewAuthHandler(authService service.IAuthService, userService service.IUserService) *AuthHandler {
	var authHandler = AuthHandler{}
	authHandler.authService = authService
	authHandler.userService = userService
	return &authHandler
}

func (h *AuthHandler) GenerateJwt(c *gin.Context) {
	// use querystring for email
	email := c.Query("email")

	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		response.ResponseNotFound(c, err.Error())
		return
	}

	token, err := h.authService.GenerateJwt(&user)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}
	response.ResponseOKWithData(c, token)
}
