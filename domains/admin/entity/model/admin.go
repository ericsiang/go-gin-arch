package model

import (
	"self_go_gin/internal/model"
)

// Admins 管理員帳號密碼表
type Admins struct {
	model.GormModel
	Account  string `gorm:"type:varchar(255);not null;uniqueIndex;" json:"account" binding:"required"`
	Password string `gorm:"type:varchar(255);not null;" json:"password" binding:"required"`
}

// func (model Admins) CreateAdmin() (Admins, error) {
// 	err := db.Create(&model).Error
// 	if err != nil {
// 		return model, err
// 	}
// 	return model, nil
// }

// func (model Admins) GetAdminsById(UserId int64) (Admins, error) {
// 	err := db.Where("id = ?", UserId).First(&model).Error
// 	if err != nil {
// 		return model, err
// 	}
// 	return model, nil
// }

// func (model Admins) GetAdminsByAccount(account string) (Admins, error) {
// 	err := db.Where("account = ?", account).First(&model).Error
// 	if err != nil {
// 		return model, err
// 	}
// 	return model, nil
// }
