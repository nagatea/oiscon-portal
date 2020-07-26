
package model

import (
	"errors"
)

// User userの構造体
type User struct {
	GormModel
	Name            string `gorm:"type:varchar(32);unique;not null;size:50" json:"name"`
	DisplayName     string `gorm:"type:varchar(64);not null" json:"displayName"`
	ProfileImageURL string `gorm:"type:varchar(64)" json:"profileImageURL"`
	Admin           bool   `gorm:"default:false" json:"admin"`
}

// RequestPutUsersBody 名前変更用のリクエストボディ
type RequestPutUsersBody struct {
	DisplayName string `json:"displayName"`
}

// TableName dbのテーブル名を指定する
func (user *User) TableName() string {
	return "users"
}

// GetUsers 全userを取得する
func GetUsers() []User {
	res := []User{}
	db.Find(&res)
	return res
}

// GetUserByName userをNameから取得する
func GetUserByName(name string) (User, error) {
	res := User{}
	db.Where("name = ?", name).First(&res)
	if name == "" {
		return User{}, errors.New("nameが存在しません")
	}
	return res, nil
}

// GetUserByID userをIDから取得する
func GetUserByID(id int) (User, error) {
	res := User{}
	db.Where("id = ?", id).First(&res)
	if res.Name == "" {
		return User{}, errors.New("該当するNameがありません")
	}
	return res, nil
}

// CreateUser userを作成する
func CreateUser(user User) (User, error) {
	if user.Name == "" {
		return User{}, errors.New("Nameが存在しません")
	}
	db.Create(&user)
	return user, nil
}

// UpdateUser userの情報(表示される名前)の変更
func UpdateUser(newUser User) (User, error) {
	user := User{}
	if newUser.Name == "" {
		return User{}, errors.New("Nameが存在しません")
	}
	db.Where("name = ?", newUser.Name).Find(&user)
	if user.Name == "" {
		return User{}, errors.New("指定したuserが存在しません")
	}
	user.DisplayName = newUser.DisplayName
	user.ProfileImageURL = newUser.ProfileImageURL
	user.Admin = newUser.Admin
	db.Save(&user)
	return user, nil
}
