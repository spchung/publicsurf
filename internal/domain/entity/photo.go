package entity

import (
	"strings"
	"time"

	"gorm.io/datatypes"
)

type Photo struct {
	ID          uint64         `gorm:"primary_key:auto_increment" json:"id" gorm:"column:id"`
	UUID        string         `gorm:"type:varchar(255)" json:"uuid" gorm:"column:uuid"`
	UserID      uint64         `gorm:"type:varchar(255)" json:"user_id" gorm:"column:user_id"`
	Name        string         `gorm:"type:varchar(255)" json:"name" gorm:"column:name"`
	S3Path      string         `gorm:"type:varchar(255)" json:"s3_path" gorm:"column:s3_path"`
	FolderID    string         `gorm:"type:varchar(255)" json:"folder_id" gorm:"column:folder_id"`
	PhotoTypeID uint64         `gorm:"type:varchar(255)" json:"photo_type_id" gorm:"column:photo_type_id"`
	Metadata    datatypes.JSON `gorm:"type:json" json:"metadata" gorm:"column:metadata"`
	PricingData datatypes.JSON `gorm:"type:json" json:"pricing_data" gorm:"column:pricing_data"`
	CreatedAt   *time.Time     `gorm:"type:timestamp" json:"created_at" gorm:"column:created_at default:null"`
	UpdatedAt   *time.Time     `gorm:"type:timestamp" json:"updated_at" gorm:"column:updated_at default:null"`
	DeletedAt   *time.Time     `gorm:"type:timestamp" json:"deleted_at" gorm:"column:deleted_at default:null"`
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
