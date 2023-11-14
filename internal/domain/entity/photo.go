package entity

import (
	"strings"

	"gorm.io/datatypes"
)

type Photo struct {
	ID          uint64         `gorm:"primary_key:auto_increment" json:"id" gorm:"column:id"`
	UserID      uint64         `gorm:"type:varchar(255)" json:"user_id" gorm:"column:user_id"`
	Name        string         `gorm:"type:varchar(255)" json:"name" gorm:"column:name"`
	S3Path      string         `gorm:"type:varchar(255)" json:"s3_path" gorm:"column:s3_path"`
	FolderID    string         `gorm:"type:varchar(255)" json:"folder_id" gorm:"column:folder_id"`
	Metadata    datatypes.JSON `gorm:"type:json" json:"metadata" gorm:"column:metadata"`
	PricingData datatypes.JSON `gorm:"type:json" json:"pricing_data" gorm:"column:pricing_data"`
}

type PhotoFolder struct {
	ID          uint64 `gorm:"primary_key:auto_increment" json:"id" gorm:"column:id"`
	Name        string `gorm:"type:varchar(255)" json:"name" gorm:"column:name"`
	Description string `gorm:"type:varchar(255)" json:"description" gorm:"column:description"`
}

type PhotoViewModel struct {
	ID     uint64 `json:"id"`
	UserID uint64 `json:"user_id"`
	Name   string `json:"name"`
	S3Path string `json:"s3_path"`
	URL    string `json:"url"`
}

func (p *PhotoViewModel) Validate(action string) map[string]string {
	errMessages := make(map[string]string)
	switch strings.ToLower(action) {
	case "create":
		break
	case "update":
		break
	default:
		break
	}
	return errMessages
}
