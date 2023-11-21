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

func (h *PhotoHandler) ProcessImage(c *gin.Context) {

	config := config.NewConfig()
	dir := config.Images.HdPath
	imageName := "water.jpg"
	_, err := h.photoService.GenerateImages(dir, imageName)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}
	response.ResponseOK(c, "success")
}

func (h *PhotoHandler) UploadPhoto(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusBadRequest)
		return
	}
	userID := c.PostForm("user_id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusBadRequest)
		return
	}
	photo, err := h.photoService.UploadPhoto(file, uint64(id))
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}
	response.ResponseOK(c, photo)
}

func (h *PhotoHandler) SaveFileToDisk(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusBadRequest)
		return
	}
	config := config.NewConfig()
	err = h.photoService.SaveFileToDisk(file, config.Images.HdPath, file.Filename)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}
	response.ResponseOK(c, "success")
}
