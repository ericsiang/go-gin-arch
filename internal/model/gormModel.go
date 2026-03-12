package model

import (
	"time"

	"gorm.io/gorm"
)

// GormModel 是所有 GORM 模型的基礎結構，包含常用的欄位
type GormModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
