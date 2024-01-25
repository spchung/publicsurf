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

	health := router.Group("/health")
	{
		health.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "OK",
			})
		})
	}

	//Register repository
	userRepo := repository.NewUserRepository(db)
	photoRepo := repository.NewPhotoRepository(db)

	router.Use(middleware.CORSMiddleware())

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
	photoAPI.Use(middleware.AuthMiddleware(db))
	photoService := service.NewPhotoService(userRepo, photoRepo)
	photoHandler := handler.NewPhotoHandler(photoService, userService)
	{
		photoAPI.GET("/list/:user_id", photoHandler.ListUserPhotos)
		photoAPI.GET("/:id", photoHandler.GetPhoto)
	}

	creatorAPI := baseAPI.Group("/creator")
	creatorAPI.Use(middleware.IsCreatorMiddleware(db))
	{
		creatorAPI.POST("/upload", photoHandler.GenerateAndUploadImages)
	}

	authAPI := baseAPI.Group("/auth")
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService, userService)
	{
		authAPI.GET("/jwt", authHandler.GenerateJwt)
	}

	return router
}
