package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"public-surf/internal/domain/repository"
	"public-surf/internal/utils"
	"public-surf/pkg/config"
)

type PhotoService struct {
	userRepo  repository.IUserRepository
	photoRepo repository.IPhotoRepository
}

type IPhotoService interface {
	GetPhotoUploaderName(id uint64) (string, error)
	GenerateImages(dir string, imageName string) error
	UploadPhoto(file *multipart.FileHeader, userID uint64) (string, error)
	SaveFileToDisk(file *multipart.FileHeader, dir string, fileName string) error
}

func NewPhotoService(userRepo repository.IUserRepository, photoRepo repository.IPhotoRepository) *PhotoService {
	return &PhotoService{userRepo: userRepo, photoRepo: photoRepo}
}

func (s *PhotoService) GetPhotoUploaderName(photoID uint64) (string, error) {
	photo, err := s.photoRepo.FindByID(photoID)
	if err != nil {
		return "", err
	}
	user, err := s.userRepo.GetUser(photo.UserID)
	if err != nil {
		return "", err
	}
	return user.FirstName + " " + user.LastName, nil
}

// generate thumbnails and regular sized photos
func (s *PhotoService) GenerateImages(dir string, imageName string) error {

	fmt.Println(dir + imageName)

	config := config.NewConfig()
	imageBytes, err := utils.LoadImg(dir, imageName)
	if err != nil {
		return err
	}
	// generate thumbnails
	thumbnailBytes, err := utils.ResizeImg(imageBytes, 165, 165)
	if err != nil {
		return err
	}
	thumbnailBytes, err = utils.WaterMark(thumbnailBytes, nil)
	if err != nil {
		return err
	}
	thumbnailName := "thumbnail_" + imageName
	err = utils.SaveImg(thumbnailBytes, config.Images.ThumbnailPath+thumbnailName)
	if err != nil {
		return err
	}
	// generate medium sized photos
	mediumBytes, err := utils.ResizeImg(imageBytes, 600, 600)
	if err != nil {
		return err
	}
	mediumBytes, err = utils.WaterMark(mediumBytes, nil)
	if err != nil {
		return err
	}
	regularName := "regular_" + imageName
	err = utils.SaveImg(mediumBytes, config.Images.RegularPath+regularName)
	if err != nil {
		return err
	}
	return nil
}

func (s *PhotoService) UploadPhoto(file *multipart.FileHeader, userID uint64) (string, error) {
	// upload photo to s3
	// save photo to db
	return "path", nil
}

func (s *PhotoService) SaveFileToDisk(file *multipart.FileHeader, dir string, fileName string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dest := dir + fileName
	dst, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// generate images
	err = s.GenerateImages(dir, fileName)
	if err != nil {
		return err
	}

	return nil
}
