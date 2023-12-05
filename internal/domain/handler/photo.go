package handler

import (
	"net/http"
	"public-surf/internal/domain/service"
	"public-surf/pkg/config"
	"public-surf/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PhotoHandler struct {
	photoService service.IPhotoService
}

func NewPhotoHandler(photoService service.IPhotoService) *PhotoHandler {
	var photoHandler = PhotoHandler{}
	photoHandler.photoService = photoService
	return &photoHandler
}

func (h *PhotoHandler) GetPhotoUploaderName(c *gin.Context) {
	photoID := c.Param("id")
	id, err := strconv.Atoi(photoID)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusBadRequest)
		return
	}
	name, err := h.photoService.GetPhotoUploaderName(uint64(id))
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}
	response.ResponseOK(c, name)
}

func (h *PhotoHandler) GenerateAndUploadImages(c *gin.Context) {

	config := config.NewConfig()
	dir := config.Images.HdPath
	imageName := "water.jpg"
	_, err := h.photoService.GenerateAndUploadImages(dir, imageName)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}
	response.ResponseOK(c, "success")
}

func (h *PhotoHandler) ListUserPhotos(c *gin.Context) {
	// get user
	userEmail, exists := c.Get("user_email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
		return
	}

	_, err := h.photoService.ListUserPhotos(userEmail.(string))
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}
	response.ResponseOK(c, "success")
}
