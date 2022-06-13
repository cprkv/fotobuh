package db

import (
	"errors"
	"time"
)

func (db *DatabaseContext) CreateCategory(name string) (uint, error) {
	if len(name) == 0 {
		return 0, errors.New("category name was empty")
	}
	cat := Category{Name: name, CreatedAt: time.Now()}
	r := db.Create(&cat)
	return cat.ID, r.Error
}

func (db *DatabaseContext) GetCategory(id uint) (Category, error) {
	cat := Category{ID: id}
	r := db.First(&cat)
	return cat, r.Error
}


func (db *DatabaseContext) GetCategories() ([]Category, error) {
	cat := []Category{}
	r := db.Find(&cat)
	return cat, r.Error
}
