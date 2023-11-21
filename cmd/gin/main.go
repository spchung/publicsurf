package main

import (
	"log"
	"public-surf/api"
	docs "public-surf/docs"
	"public-surf/pkg/config"
	"public-surf/pkg/database"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	config.GetConfig()
}

// @BasePath /api/v1
func main() {
	db, err := database.InitDB()
	docs.SwaggerInfo.BasePath = "/api/v1"
	if err != nil {
		log.Fatal(err)
	}
	r := api.SetupRouter(db)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")
}
