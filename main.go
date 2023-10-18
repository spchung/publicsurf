package main

import (
	"log"
	"public-surf/api"
	"public-surf/pkg/config"
	"public-surf/pkg/database"
)

func init() {
	config.GetConfig()
}

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	r := api.SetupRouter(db)
	r.Run(":8080")
}
