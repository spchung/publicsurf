package handler

import (
	"net/http"
	"public-surf/internal/domain/entity"
	"public-surf/internal/domain/service"
	"public-surf/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PhotoHandler struct {
	photoService service.IPhotoService
}

type UploadPhotoRequest struct {
	Name string `json:"name"`
}

func NewPhotoHandler(photoService service.IPhotoService) *PhotoHandler {
	var photoHandler = PhotoHandler{}
	photoHandler.photoService = photoService
	return &photoHandler
}

func (h *PhotoHandler) GenerateAndUploadImages(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusBadRequest)
		return
	}
	imageName := c.PostForm("name")

	savesPhoto, err := h.photoService.GenerateAndUploadImages(file, imageName)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}
	response.ResponseOKWithData(c, savesPhoto)
}

func (h *PhotoHandler) ListUserPhotos(c *gin.Context) {
	userEmail, exists := c.Get("user_email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
		return
	}

	photos, err := h.photoService.ListUserPhotos(userEmail.(string))
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}
	res := make([]entity.PhotoViewModel, len(photos))
	for i, photo := range photos {
		res[i] = entity.PhotoViewModel{
			ID:     photo.ID,
			UserID: photo.UserID,
			Name:   photo.Name,
			S3Path: photo.S3Path,
		}
	}

	response.ResponseOKWithData(c, res)
}

func (h *PhotoHandler) GetPhoto(c *gin.Context) {
	photoID := c.Param("id")
	// convert string to uint64
	photoIDUint, err := strconv.ParseUint(photoID, 10, 64)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}
	photo, err := h.photoService.GetPhoto(photoIDUint)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}
	response.ResponseOKWithData(c, photo)
}
