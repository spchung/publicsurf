package handler

import (
	"net/http"
	"public-surf/internal/domain/service"
	"public-surf/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.IUserService
}

type IUserHandler interface {
	GetUserPhotoCount(c *gin.Context)
	GetUserLatestPhoto(c *gin.Context)
	GetUser(c *gin.Context)
}

func NewUserHandler(userService service.IUserService) *UserHandler {
	var userHandler = UserHandler{}
	userHandler.userService = userService
	return &userHandler
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := h.userService.GetUser(id)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}
	response.ResponseOKWithData(c, user)
}

func (h *UserHandler) GetUserPhotoCount(c *gin.Context) {
	userID := c.Param("id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusBadRequest)
		return
	}
	count, err := h.userService.GetUserPhotoCount(id)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}
	response.ResponseOKWithData(c, count)
}

func (h *UserHandler) GetUserLatestPhoto(c *gin.Context) {
	userID := c.Param("id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusBadRequest)
		return
	}
	photo, err := h.userService.GetUserLatestPhoto(id)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}
	response.ResponseOKWithData(c, photo)
}
