package model

import "gorm.io/gorm"

var db *gorm.DB

// SetDB 設定 GORM 資料庫連接
func SetDB(gorm *gorm.DB) {
	db = gorm
}
