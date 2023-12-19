package main

import (
	"log"
	"public-surf/api"
	"public-surf/internal/logger"
	"public-surf/pkg/config"
	"public-surf/pkg/database"

	"go.uber.org/zap"
)

func init() {
	config.GetConfig()
}

// @BasePath /api/v1
func main() {
	logger := logger.NewLogger()
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	r := api.SetupRouter(db)
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// docs.SwaggerInfo.BasePath = "/api/v1"
	r.Run("0.0.0.0:8080")
}
