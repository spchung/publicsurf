package entity

type Cart struct {
	ID          int    `gorm:"primary_key:auto_increment" json:"id" gorm:"column:id"`
	UserID      int    `gorm:"type:varchar(255)" json:"user_id" gorm:"column:user_id"`
	Description string `gorm:"type:varchar(255)" json:"description" gorm:"column:description"`
}

type CartPhoto struct {
	ID      int `gorm:"primary_key:auto_increment" json:"id" gorm:"column:id"`
	CartID  int `gorm:"type:varchar(255)" json:"cart_id" gorm:"column:cart_id"`
	PhotoID int `gorm:"type:varchar(255)" json:"photo_id" gorm:"column:photo_id"`
}
