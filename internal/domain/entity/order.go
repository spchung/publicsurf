package entity

type Order struct {
	ID          uint64 `gorm:"primary_key:auto_increment" json:"id" gorm:"column:id"`
	UserID      uint64 `gorm:"type:varchar(255)" json:"user_id" gorm:"column:user_id"`
	Description string `gorm:"type:varchar(255)" json:"description" gorm:"column:description"`
}

type OrderPhoto struct {
	ID      uint64 `gorm:"primary_key:auto_increment" json:"id" gorm:"column:id"`
	OrderID uint64 `gorm:"type:varchar(255)" json:"order_id" gorm:"column:order_id"`
	PhotoID uint64 `gorm:"type:varchar(255)" json:"photo_id" gorm:"column:photo_id"`
}
