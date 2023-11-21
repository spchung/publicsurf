package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"public-surf/internal/domain/entity"
	"public-surf/internal/domain/repository"
	"public-surf/internal/utils"
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
	GenerateImages(dir string, imageName string) (*entity.Photo, error)
	UploadPhoto(file *multipart.FileHeader, userID uint64) (string, error)
	SaveFileToDisk(file *multipart.FileHeader, dir string, fileName string) error
}

func NewPhotoService(userRepo repository.IUserRepository, photoRepo repository.IPhotoRepository) *PhotoService {
	return &PhotoService{
		userRepo:  userRepo,
		photoRepo: photoRepo,
	}
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

func (s *PhotoService) uploadImage(imageBytes []byte, photoUuid string, imageName string) error {
	config := config.NewConfig()

	s3Client := aws_helper.NewAwsS3Client()

	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(config.Files.PublicBucket),
		Key:    aws.String(fmt.Sprintf("images/%s/%s", photoUuid, imageName)),
		Body:   bytes.NewReader(imageBytes),
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// generate thumbnails and regular sized photos
func (s *PhotoService) GenerateImages(dir string, imageName string) (*entity.Photo, error) {

	// init a uuid for images to use as db photo uuid and s3 direcotry name
	uuid := uuid.New().String()
	wg := sync.WaitGroup{}
	imageBytes, err := utils.LoadImg(dir, imageName)
	if err != nil {
		return nil, err
	}

	wg.Add(2)
	go func() error {
		defer wg.Done()

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
		err = s.uploadImage(thumbnailBytes, uuid, thumbnailName)
		if err != nil {
			return err
		}
		return nil
	}()

	go func() error {
		defer wg.Done()
		mediumBytes, err := utils.ResizeImg(imageBytes, 600, 600)
		if err != nil {
			return err
		}
		mediumBytes, err = utils.WaterMark(mediumBytes, nil)
		if err != nil {
			return err
		}
		regularName := "regular_" + imageName
		err = s.uploadImage(mediumBytes, uuid, regularName)
		if err != nil {
			return err
		}
		return nil
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
	_, err = s.GenerateImages(dir, fileName)
	if err != nil {
		return err
	}

	return nil
}
