package productdomain

import (
	"ecommerce/common"
)

type Product struct {
	common.BaseModel
	CategoryID  int     `gorm:"column:category_id" json:"category_id"`
	Name        string  `gorm:"column:name" json:"name"`
	Image       *string `gorm:"column:image" json:"image"`
	Type        string  `gorm:"column:type" json:"type"`
	Description string  `gorm:"column:description" json:"description"`
}