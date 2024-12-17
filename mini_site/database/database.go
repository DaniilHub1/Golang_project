package database

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "mini_site/models"
)

var DB *gorm.DB

// Инициализация базы данных и миграции
func InitDB() (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    DB = db

    // Автоматическая миграция для создания таблицы
    err = db.AutoMigrate(&models.Post{})
    if err != nil {
        return nil, err
    }

    return DB, nil
}
