package api

import (
	"public-surf/api/middleware"
	"public-surf/internal/domain/handler"
	"public-surf/internal/domain/repository"
	"public-surf/internal/domain/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {

	router := gin.Default()
	//Register repository
	userRepo := repository.NewUserRepository(db)
	photoRepo := repository.NewPhotoRepository(db)

	baseAPI := router.Group("/v1")
	userAPI := baseAPI.Group("/user")
	userService := service.NewUserService(userRepo, photoRepo)
	userHandler := handler.NewUserHandler(userService)
	{
		userAPI.GET("/:id", userHandler.GetUser)
		userAPI.GET("/:id/count", userHandler.GetUserPhotoCount)
		userAPI.GET("/:id/latest-photo", userHandler.GetUserLatestPhoto)
	}

	photoAPI := baseAPI.Group("/photo")
	photoAPI.Use(middleware.AuthMiddleware())
	photoService := service.NewPhotoService(userRepo, photoRepo)
	photoHandler := handler.NewPhotoHandler(photoService)
	{
		// photoAPI.POST("/upload", photoHandler.UploadPhoto)
		photoAPI.GET("/list", photoHandler.ListUserPhotos)
		photoAPI.GET("/:id/uploader-name", photoHandler.GetPhotoUploaderName)
		photoAPI.GET("/generate_upload", photoHandler.GenerateAndUploadImages)
	}

	authAPI := baseAPI.Group("/auth")
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService, userService)
	{
		authAPI.GET("/jwt", authHandler.GenerateJwt)
	}

	return router
}
