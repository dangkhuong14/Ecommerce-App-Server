package domain

import (

)

type ProductUpdatingDTO struct {
	CategoryID  *int    `gorm:"column:category_id" json:"category_id"`
	Name        *string `gorm:"column:name" json:"name"`
	Image       *string `gorm:"column:image" json:"image"`
	Type        *string `gorm:"column:type" json:"type"`
	Description *string `gorm:"column:description" json:"description"`
	Status      *string `gorm:"column:status" json:"status"`
}

type ProductCreationDTO struct {
	CategoryID  int    `gorm:"column:category_id" json:"category_id"`
	Name        string `gorm:"column:name" json:"name"`
	Image       string `gorm:"column:image" json:"image"`
	Type        string `gorm:"column:type" json:"type"`
	Description string `gorm:"column:description" json:"description"`
}