package db

import (
	"errors"
	"fotobuh/lib"
	"os"
	"path/filepath"
	"time"

	"github.com/jbrodriguez/mlog"
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

func (db *DatabaseContext) GetCategoryWithPictures(id uint) (Category, error) {
	cat := Category{ID: id}
	r := db.Preload("Pictures").First(&cat)
	return cat, r.Error
}

func (db *DatabaseContext) GetCategories() ([]Category, error) {
	cat := []Category{}
	r := db.Find(&cat)
	return cat, r.Error
}

func (db *DatabaseContext) DeleteCategory(c *Category) error {
	catWithPics, err := db.GetCategoryWithPictures(c.ID)
	if err != nil {
		return err
	}

	picIds := make([]uint, len(catWithPics.Pictures))
	for i := 0; i < len(catWithPics.Pictures); i++ {
		picIds[i] = catWithPics.Pictures[i].ID
		path := filepath.Join(lib.Config.Pictures.Storage, catWithPics.Pictures[i].FileName)
		mlog.Trace("removing %v", path)
		err = os.Remove(path)
		if err != nil {
			mlog.Warning("erro removing file '%v', remove it manually!", path)
		}
	}

	r := db.Delete(&Picture{}, picIds)
	if r.Error != nil {
		return r.Error
	}

	r = db.Delete(c)
	return r.Error
}

func (db *DatabaseContext) CreatePicture(p *Picture) error {
	r := db.Create(p)
	return r.Error
}
