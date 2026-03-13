// Package mysqlmgr 提供了 MySQL 錯誤處理與 GORM 記錄檢查的工具函數
package mysqlmgr

import (
	mysqlErr "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// DuplicateEntryCode 重複條目錯誤代碼
const DuplicateEntryCode = 1062

// MysqlErrCode 根据mysql错误信息返回错误代码
/*
* 1062: Duplicate entry
 */
func MysqlErrCode(err error) int {
	mysqlErr, ok := err.(*mysqlErr.MySQLError)
	if !ok {
		return 0
	}
	return int(mysqlErr.Number)
}

// CheckRecordNotFound 檢查是否沒有找到記錄
func CheckRecordNotFound(result *gorm.DB) error {
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
