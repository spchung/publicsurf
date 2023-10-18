package database

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	username := viper.GetString("Database.Username")
	password := viper.GetString("Database.Password")
	host := viper.GetString("Database.Host")
	port := viper.GetInt("Database.Port")
	dbname := viper.GetString("Database.DBName")

	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, dbname)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		host,
		username,
		password,
		dbname,
		port,
	)
	// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	return db, nil

}
