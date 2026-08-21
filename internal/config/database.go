package config

import (
	"github.com/erfangho/url-shortener/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("url_shortener.db"), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.URL{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
