package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"public-surf/internal/domain/entity"
	"public-surf/internal/domain/repository"
	"public-surf/internal/utils"
	"strings"
	"sync"
	"time"

	"public-surf/pkg/aws_helper"
	"public-surf/pkg/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type PhotoService struct {
	userRepo  repository.IUserRepository
	photoRepo repository.IPhotoRepository
}

type IPhotoService interface {
	GetPhotoUploaderName(id uint64) (string, error)
	ListUserPhotos(userEmail string) ([]*entity.Photo, error)
	GenerateAndUploadImages(file *multipart.FileHeader, imageName string) ([]*entity.Photo, error)
	GetPhoto(id uint64) (*entity.PhotoView, error)
}

func NewPhotoService(userRepo repository.IUserRepository, photoRepo repository.IPhotoRepository) *PhotoService {
	return &PhotoService{
		userRepo:  userRepo,
		photoRepo: photoRepo,
	}
}

func (s *PhotoService) ListUserPhotos(userEmail string) ([]*entity.Photo, error) {
	// get user
	user, err := s.userRepo.GetUserByEmail(userEmail)
	if err != nil {
		return nil, err
	}

	photos, err := s.photoRepo.FindByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	fmt.Println(userEmail)

	return photos, nil
}

func (s *PhotoService) GetPhotoUploaderName(photoID uint64) (string, error) {
	photo, err := s.photoRepo.GetByID(photoID)
	if err != nil {
		return "", err
	}
	user, err := s.userRepo.GetUser(photo.UserID)
	if err != nil {
		return "", err
	}
	return user.FirstName + " " + user.LastName, nil
}

func (s *PhotoService) GenerateAndUploadImages(file *multipart.FileHeader, imageName string) ([]*entity.Photo, error) {
	// uploadCompleteChan := make(chan *entity.Photo, 3)
	errorChan := make(chan error, 3)

	// init a uuid for images to use as db photo uuid and s3 direcotry name
	uuid := uuid.New().String()
	wg := sync.WaitGroup{}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	imageBytes, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	wg.Add(3)
	now := time.Now()
	successfulPhotos := []*entity.Photo{}
	// thumbnail
	go func() {
		defer wg.Done()
		thumbnailBytes, err := s.generateThumbnail(imageBytes)
		if err != nil {
			errorChan <- err
			return
		}

		thumbnailName := "thumbnail_" + imageName
		s3_path, err := s.uploadImagePublic(thumbnailBytes, uuid, thumbnailName)
		if err != nil {
			errorChan <- err
			return
		}
		savedPhoto, err := s.photoRepo.Save(&entity.Photo{
			UUID:        uuid,
			Name:        imageName,
			UserID:      1,
			PhotoTypeID: 1,
			S3Path:      s3_path,
			CreatedAt:   &now,
			UpdatedAt:   nil,
		})
		if err != nil {
			errorChan <- err
			return
		}
		successfulPhotos = append(successfulPhotos, savedPhoto)
	}()

	// regular sized image
	go func() {
		defer wg.Done()
		mediumBytes, err := s.generateregular(imageBytes)
		if err != nil {
			errorChan <- err
			return
		}
		regularName := "regular_" + imageName
		s3_path, err := s.uploadImagePublic(mediumBytes, uuid, regularName)
		if err != nil {
			errorChan <- err
			return
		}
		savedPhoto, err := s.photoRepo.Save(&entity.Photo{
			UUID:        uuid,
			Name:        imageName,
			UserID:      1,
			PhotoTypeID: 2,
			S3Path:      s3_path,
			CreatedAt:   &now,
			UpdatedAt:   nil,
		})
		if err != nil {
			errorChan <- err
			return
		}
		successfulPhotos = append(successfulPhotos, savedPhoto)

	}()

	// original image - private bucket
	go func() {
		defer wg.Done()
		s3_path, err := s.uploadImagePivate(imageBytes, uuid, imageName)
		if err != nil {
			errorChan <- err
			return
		}
		savedPhoto, err := s.photoRepo.Save(&entity.Photo{
			UUID:        uuid,
			Name:        imageName,
			UserID:      1,
			PhotoTypeID: 3,
			S3Path:      s3_path,
			CreatedAt:   &now,
			UpdatedAt:   nil,
		})
		if err != nil {
			errorChan <- err
			return
		}
		successfulPhotos = append(successfulPhotos, savedPhoto)
	}()

	// listen for errors and upload complete
	// errorList := []string{}

	go func() {
		for {
			select {
			case err, ok := <-errorChan:
				if !ok {
					return
				}
				fmt.Println(err)
			}
		}
	}()

	// TODO: make sure no upload errors
	wg.Wait()

	return successfulPhotos, nil
}

func (s *PhotoService) GetPhoto(id uint64) (*entity.PhotoView, error) {
	photo, err := s.photoRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return &photo, nil
}

// uploads to the public bucket - "public-surf"
func (s *PhotoService) uploadImagePublic(imageBytes []byte, photoUuid string, imageName string) (string, error) {
	config := config.NewConfig()

	return s.uploadImage(config.Files.PublicBucket, imageBytes, photoUuid, imageName)
}

func (s *PhotoService) uploadImagePivate(imageBytes []byte, photoUuid string, imageName string) (string, error) {
	config := config.NewConfig()
	return s.uploadImage(config.Files.PrivateBucket, imageBytes, photoUuid, imageName)
}

func (s *PhotoService) uploadImage(bucket string, imageBytes []byte, photoUuid string, imageName string) (string, error) {
	s3Client := aws_helper.NewAwsS3Client()

	dir := fmt.Sprintf("images/%s/%s", photoUuid, imageName)
	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(dir),
		Body:   bytes.NewReader(imageBytes),
	})
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return strings.Join([]string{bucket, dir}, "/"), nil
}

// sub routines for process different image sizes
func (s *PhotoService) generateThumbnail(imageBytes []byte) ([]byte, error) {
	thumbnailBytes, err := utils.ResizeImg(imageBytes, 165, 165)
	if err != nil {
		return nil, err
	}
	thumbnailBytes, err = utils.WaterMark(thumbnailBytes, nil)
	if err != nil {
		return nil, err
	}
	return thumbnailBytes, nil
}

func (s *PhotoService) generateregular(imageBytes []byte) ([]byte, error) {
	mediumBytes, err := utils.ResizeImg(imageBytes, 600, 600)
	if err != nil {
		return nil, err
	}
	mediumBytes, err = utils.WaterMark(mediumBytes, nil)
	if err != nil {
		return nil, err
	}
	return mediumBytes, nil
}
