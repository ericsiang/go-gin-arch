// Package dao 定義持久化物件（Persistence Object）
package dao

import (
	"self_go_gin/internal/model"
)

// UserPO 用戶持久化物件
// 用途：僅用於資料庫 ORM 映射，不包含業務邏輯
// 與領域模型（User 聚合根）分離，避免領域層被基礎設施層污染
type UserPO struct {
	model.GormModel
	Account  string `gorm:"type:varchar(255);not null;uniqueIndex;column:account" json:"account"`
	Password string `gorm:"type:varchar(255);not null;column:password" json:"password"`
}

// TableName 指定表名
func (UserPO) TableName() string {
	return "users"
}
