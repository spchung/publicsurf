package entity

type Order struct {
	ID          int    `gorm:"primary_key:auto_increment" json:"id" gorm:"column:id"`
	UserID      int    `gorm:"type:varchar(255)" json:"user_id" gorm:"column:user_id"`
	Description string `gorm:"type:varchar(255)" json:"description" gorm:"column:description"`
}

type OrderPhoto struct {
	ID      int `gorm:"primary_key:auto_increment" json:"id" gorm:"column:id"`
	OrderID int `gorm:"type:varchar(255)" json:"order_id" gorm:"column:order_id"`
	PhotoID int `gorm:"type:varchar(255)" json:"photo_id" gorm:"column:photo_id"`
}
