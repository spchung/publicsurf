package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig
	Images   ImagesConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type ImagesConfig struct {
	BasePath      string
	HdPath        string
	ThumbnailPath string
	RegularPath   string
}

func GetConfig() {
	viper.SetConfigName("App")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/Users/stephen/Dev/public_surf/backend/configurations")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("config error : ", err.Error())
	}
}

func NewConfig() *Config {
	GetConfig()
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("config error : ", err.Error())
	}
	return &config
}
