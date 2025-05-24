package database

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "mini_site/models"
)

var DB *gorm.DB


func InitDB() (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open("database/mini_site.db"), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    DB = db

    err = db.AutoMigrate(&models.Post{})
    if err != nil {
        return nil, err
    }

    return DB, nil
}

