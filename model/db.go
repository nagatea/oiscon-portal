package model

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

var allTables = []interface{}{
	User{},
}

// GormModel Gormが必要とする構造体
type GormModel struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

// EstablishConnection DBに接続する
func EstablishConnection() (*gorm.DB, error) {
	user := os.Getenv("MYSQL_USERNAME")
	if user == "" {
		user = "root"
	}

	pass := os.Getenv("MYSQL_PASSWORD")
	if pass == "" {
		pass = "password"
	}

	host := os.Getenv("MYSQL_HOST")
	if host == "" {
		host = "localhost"
	}

	dbname := os.Getenv("MYSQL_DATABASE")
	if dbname == "" {
		dbname = "portal"
	}

	_db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, pass, host, dbname)+"?parseTime=true&loc=Asia%2FTokyo&charset=utf8mb4")
	db = _db
	return db, err
}

// Migrate DBのマイグレーション
func Migrate() error {
	if err := db.AutoMigrate(allTables...).Error; err != nil {
		return err
	}

	admin, _ := GetUserByName("nagatea")
	if admin.Name == "" {
		user := User{
			Name:        "nagatea",
			DisplayName: "ながてち",
			Admin:       true,
		}
		_, err := CreateUser(user)
		if err != nil {
			return err
		}
	}

	return nil
}
