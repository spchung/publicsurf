package entity

import "strings"

type Photo struct {
	ID     uint64 `gorm:"primary_key:auto_increment" json:"id" gorm:"column:id"`
	UserID uint64 `gorm:"type:varchar(255)" json:"user_id" gorm:"column:user_id"`
	Name   string `gorm:"type:varchar(255)" json:"name" gorm:"column:name"`
	S3Path string `gorm:"type:varchar(255)" json:"s3_path" gorm:"column:s3_path"`
	URL    string `gorm:"type:varchar(255)" json:"url" gorm:"column:url"`
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

// func (p *Photo) ToViewModel() *PhotoViewModel {
// 	return &PhotoViewModel{
// 		ID:     p.ID,
// 		UserID: p.UserID,
// 		Name:   p.Name,
// 		S3Path: p.S3Path,
// 		URL:    p.URL,
// 	}
// }
