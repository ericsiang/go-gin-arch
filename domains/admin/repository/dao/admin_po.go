// Package dao 定義持久化物件（Persistence Object）
package dao

import (
	"self_go_gin/internal/model"
)

// AdminPO 管理員持久化物件
type AdminPO struct {
	model.GormModel
	Account  string `gorm:"type:varchar(255);not null;uniqueIndex;column:account" json:"account"`
	Password string `gorm:"type:varchar(255);not null;column:password" json:"password"`
}

// TableName 指定表名
func (AdminPO) TableName() string {
	return "admins"
}
