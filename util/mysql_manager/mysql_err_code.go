package mysql_manager

import (
	"strconv"

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

// StringToUint64 將字符串轉換為uint64
func StringToUint64(s string) (uint64, error) {
	ui64, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return ui64, nil
}

// StringToUint16 將字符串轉換為uint16
func StringToUint16(s string) (uint16, error) {
	ui64, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	ui16 := uint16(ui64)
	return ui16, nil
}
