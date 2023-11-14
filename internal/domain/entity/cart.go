package entity

type Cart struct {
	ID          uint64 `gorm:"primary_key:auto_increment" json:"id" gorm:"column:id"`
	UserID      uint64 `gorm:"type:varchar(255)" json:"user_id" gorm:"column:user_id"`
	Description string `gorm:"type:varchar(255)" json:"description" gorm:"column:description"`
}

type CartPhoto struct {
	ID      uint64 `gorm:"primary_key:auto_increment" json:"id" gorm:"column:id"`
	CartID  uint64 `gorm:"type:varchar(255)" json:"cart_id" gorm:"column:cart_id"`
	PhotoID uint64 `gorm:"type:varchar(255)" json:"photo_id" gorm:"column:photo_id"`
}
