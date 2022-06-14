package db

import (
	"time"
)

type Picture struct {
	ID         uint `gorm:"primaryKey;autoIncrement"`
	Name       string
	FileName   string
	CreatedAt  time.Time   `gorm:"index:,sort:desc"`
	Categories []*Category `gorm:"many2many:picture_categories;"`
}

type Category struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	Name      string
	CreatedAt time.Time  `gorm:"index:,sort:desc"`
	Pictures  []*Picture `gorm:"many2many:picture_categories;"`
}
