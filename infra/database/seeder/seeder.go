package seeder

import (
	"gorm.io/gorm"
)

// Seeder 數據庫種子器
type Seeder struct {
	db *gorm.DB
}

// RunSeeder 執行數據庫種子建立資料
func RunSeeder() {
	//common_seeder
	CreateUser()
	CreateAdmin()
}

// NewSeeder 創建新的數據庫種子器
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

// Clear 清空指定表格的所有資料
func (s *Seeder) Clear(tableName string) error {
	return s.db.Exec("truncate table " + tableName).Error
}
