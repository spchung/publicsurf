package service

import (
	"bytes"
	"context"
	"fmt"
	"public-surf/internal/domain/entity"
	"public-surf/internal/domain/repository"
	"public-surf/internal/utils"
	"strings"
	"sync"

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
	// GenerateImages(dir string, imageName string) (*entity.Photo, error)
	// UploadPhoto(file *multipart.FileHeader, userID uint64) (string, error)
	ListUserPhotos(userEmail string) ([]*entity.Photo, error)
	GenerateAndUploadImages(dir string, imageName string) (*entity.Photo, error)
}

func NewPhotoService(userRepo repository.IUserRepository, photoRepo repository.IPhotoRepository) *PhotoService {
	return &PhotoService{
		userRepo:  userRepo,
		photoRepo: photoRepo,
	}
}

func (s *PhotoService) ListUserPhotos(userEmail string) ([]*entity.Photo, error) {
	// get user
	// user, err := s.userRepo.GetUser(uint64(userID))
	// if err != nil {
	// 	return nil, err
	// }

	fmt.Println(userEmail)

	return nil, nil
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

func (s *PhotoService) generateMedium(imageBytes []byte) ([]byte, error) {
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

// generate thumbnails and regular sized images and upload with original
func (s *PhotoService) GenerateAndUploadImages(dir string, imageName string) (*entity.Photo, error) {

	uploadCompleteChan := make(chan string, 3)
	errorChan := make(chan error, 3)

	// init a uuid for images to use as db photo uuid and s3 direcotry name
	uuid := uuid.New().String()
	wg := sync.WaitGroup{}
	imageBytes, err := utils.LoadImg(dir, imageName)
	if err != nil {
		return nil, err
	}

	wg.Add(3)
	// thumbnail
	go func() {
		defer wg.Done()

		// generate thumbnails
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
		uploadCompleteChan <- s3_path
	}()

	// medium sized image
	go func() {
		defer wg.Done()
		mediumBytes, err := s.generateMedium(imageBytes)
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
		uploadCompleteChan <- s3_path
	}()

	// original image - private bucket
	go func() {
		defer wg.Done()
		s3_path, err := s.uploadImagePivate(imageBytes, uuid, imageName)
		if err != nil {
			errorChan <- err
			return
		}
		uploadCompleteChan <- s3_path
	}()

	// listen for errors and upload complete
	errorList := []string{}
	successUrls := []string{}
	go func() {
		for {
			select {
			case err, ok := <-errorChan:
				if !ok {
					return
				}
				// add to error logs
				errorList = append(errorList, err.Error())
			case res, ok := <-uploadCompleteChan:
				if !ok {
					return
				}
				// add to database
				successUrls = append(successUrls, res)
			}
		}
	}()

	// TODO: make sure no upload errors
	wg.Wait()

	// save to database
	photo := repository.NewPhoto()
	photo.UUID = uuid
	photo.Name = imageName
	photo.UserID = 1
	savedPhoto, err := s.photoRepo.Save(photo)

	return savedPhoto, nil
}
