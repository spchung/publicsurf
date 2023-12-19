package handler

import (
	"net/http"
	"public-surf/internal/domain/entity"
	"public-surf/internal/domain/service"
	"public-surf/internal/logger"
	"public-surf/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PhotoHandler struct {
	photoService service.IPhotoService
	userService  service.IUserService
}

type UploadPhotoRequest struct {
	Name string `json:"name"`
}

func NewPhotoHandler(photoService service.IPhotoService, userService service.IUserService) *PhotoHandler {
	var photoHandler = PhotoHandler{}
	photoHandler.photoService = photoService
	photoHandler.userService = userService
	return &photoHandler
}

func (h *PhotoHandler) GenerateAndUploadImages(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		logger.Logger.Error("error - GenerateAndUploadImages handler", zap.Error(err))
		response.ResponseError(c, err.Error(), http.StatusBadRequest)
		return
	}
	imageName := c.PostForm("name")

	savesPhoto, err := h.photoService.GenerateAndUploadImages(file, imageName)
	if err != nil {
		logger.Logger.Error("error - GenerateAndUploadImages handler", zap.Error(err))
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}
	response.ResponseOKWithData(c, savesPhoto)
}

func (h *PhotoHandler) ListUserPhotos(c *gin.Context) {
	param := c.Param("user_id")

	userID, err := strconv.Atoi(param)
	if err != nil {
		logger.Logger.Error("error - ListUserPhotos handler", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	_, err = h.userService.GetUser(userID)
	if err != nil {
		logger.Logger.Error("error - ListUserPhotos handler", zap.Error(err))
		response.ResponseNotFound(c, "User not found")
		return
	}

	photos, err := h.photoService.ListUserPhotos(userID)
	if err != nil {
		logger.Logger.Error("error - ListUserPhotos handler", zap.Error(err))
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
	photoIDUint, err := strconv.Atoi(photoID)
	if err != nil {
		logger.Logger.Error("error - GetPhoto handler", zap.Error(err))
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}
	photo, err := h.photoService.GetPhoto(photoIDUint)
	if err != nil {
		logger.Logger.Error("error - GetPhoto handler", zap.Error(err))
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}
	response.ResponseOKWithData(c, photo)
}
