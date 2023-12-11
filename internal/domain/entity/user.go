package entity

import (
	"public-surf/pkg/security"

	"github.com/badoux/checkmail"
	"gorm.io/datatypes"
)

type User struct {
	ID           uint64         `gorm:"primary_key:auto_increment" json:"id" gorm:"column:beast_id"`
	Email        string         `gorm:"type:varchar(255)" json:"email" gorm:"column:email"`
	Password     string         `gorm:"type:varchar(255)" json:"password" gorm:"column:password"`
	FirstName    string         `gorm:"type:varchar(255)" json:"first_name" gorm:"column:first_name"`
	LastName     string         `gorm:"type:varchar(255)" json:"last_name" gorm:"column:last_name"`
	UserTypeID   uint64         `gorm:"type:int4" gorm:"foreignKey:UserTypeID" json:"user_type_id" gorm:"column:user_type_id"`
	ThumbnailUrl string         `gorm:"type:varchar(255)" json:"thumbnail_url" gorm:"column:thumbnail_url"`
	PaymentInfo  datatypes.JSON `gorm:"type:json" json:"payment_info" gorm:"column:payment_info"`
}

type UserType struct {
	ID   uint64 `gorm:"primary_key:auto_increment" json:"id" gorm:"column:id"`
	Name string `gorm:"type:varchar(255)" json:"name" gorm:"column:name"`
}

type ReqisterViewModel struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserLoginViewModel struct {
	Email    string `gorm:"type:varchar(255)" json:"email"`
	Password string `gorm:"type:varchar(255)" json:"password"`
}

type UserViewModel struct {
	ID        uint64 `gorm:"primary_key:auto_increment" json:"id"`
	Email     string `gorm:"type:varchar(255)" json:"email"`
	FirstName string `gorm:"type:varchar(255)" json:"first_name"`
	LastName  string `gorm:"type:varchar(255)" json:"last_name"`
}

func (u *User) EncryptPassword(password string) (string, error) {
	hashPassword, err := security.Hash(password)
	if err != nil {
		return "", err
	}

	return string(hashPassword), nil
}

func (u *User) Validate() map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	if u.Email == "" {
		errorMessages["email_required"] = "email required"
	}
	if u.Email != "" {
		if err = checkmail.ValidateFormat(u.Email); err != nil {
			errorMessages["invalid_email"] = "email email"
		}
	}

	return errorMessages
}

func (u *ReqisterViewModel) Validate() map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	if u.FirstName == "" {
		errorMessages["firstname_required"] = "first name is required"
	}
	if u.LastName == "" {
		errorMessages["lastname_required"] = "last name is required"
	}
	if u.Password == "" {
		errorMessages["password_required"] = "password is required"
	}
	if u.Password != "" && len(u.Password) < 6 {
		errorMessages["invalid_password"] = "password should be at least 6 characters"
	}
	if u.Email == "" {
		errorMessages["email_required"] = "email is required"
	}
	if u.Email != "" {
		if err = checkmail.ValidateFormat(u.Email); err != nil {
			errorMessages["invalid_email"] = "please provide a valid email"
		}
	}

	return errorMessages
}
