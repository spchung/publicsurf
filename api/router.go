package api

import (
	"public-surf/internal/domain/handler"
	"public-surf/internal/domain/repository"
	"public-surf/internal/domain/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {

	router := gin.Default()

	//Register User Repo
	userRepo := repository.NewUserRepository(db)
	photoRepo := repository.NewPhotoRepository(db)

	userService := service.NewUserService(userRepo, photoRepo)
	userHandler := handler.NewUserHandler(userService)

	// Register Photo Repo
	photoService := service.NewPhotoService(userRepo, photoRepo)
	photoHandler := handler.NewPhotoHandler(photoService)

	baseAPI := router.Group("/v1")
	userAPI := baseAPI.Group("/user")
	{
		userAPI.GET("/:id", userHandler.GetUser)
		userAPI.GET("/:id/count", userHandler.GetUserPhotoCount)
		userAPI.GET("/:id/latest-photo", userHandler.GetUserLatestPhoto)
	}

	photoAPI := baseAPI.Group("/photo")
	{
		photoAPI.POST("/upload", photoHandler.SaveFileToDisk)
		photoAPI.GET("/:id/uploader-name", photoHandler.GetPhotoUploaderName)
		photoAPI.GET("/process", photoHandler.ProcessImage)
	}

	return router
}
