package handler

import (
	"net/http"
	"public-surf/internal/domain/service"
	"public-surf/pkg/config"
	"public-surf/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PhotoHnadler struct {
	photoService service.IPhotoService
}

func NewPhotoHandler(photoService service.IPhotoService) *PhotoHnadler {
	var photoHandler = PhotoHnadler{}
	photoHandler.photoService = photoService
	return &photoHandler
}

func (h *PhotoHnadler) GetPhotoUploaderName(c *gin.Context) {
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

func (h *PhotoHnadler) ProcessImage(c *gin.Context) {

	config := config.NewConfig()
	dir := config.Images.HdPath
	imageName := "water.jpg"
	err := h.photoService.GenerateImages(dir, imageName)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}
	response.ResponseOK(c, "success")
}

func (h *PhotoHnadler) UploadPhoto(c *gin.Context) {
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

func (h *PhotoHnadler) SaveFileToDisk(c *gin.Context) {
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
