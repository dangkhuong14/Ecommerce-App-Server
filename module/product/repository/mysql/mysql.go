package mysql

import (
	"gorm.io/gorm"
)

type MysqlRepository struct {
	db *gorm.DB
}

func NewMysqlRepository(db *gorm.DB) MysqlRepository {
	return MysqlRepository{
		db: db,
	}
}