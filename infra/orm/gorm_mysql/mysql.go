// Package gormysql 提供了初始化 GORM 連接 MySQL 資料庫的功能
package gormysql

import (
	"fmt"
	constset "self_go_gin/common/const"
	"self_go_gin/infra/database"
	"self_go_gin/infra/env"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

// MysqlDB 是 Database 接口的實現，封裝了 GORM 的 MySQL 連接
type MysqlDB struct {
	db *gorm.DB
}

// InitMysqlDB 初始化 MysqlDB
func InitMysqlDB() database.Database {
	return &MysqlDB{}
}

// Connect 建立 MySQL 連接
func (m *MysqlDB) Connect(serverEnv *env.ServerConfig) error {
	var config *gorm.Config
	gormZaplogger := zapgorm2.New(zap.L())
	logger.Default.LogMode(logger.Error)
	// zap.S().Info("logger level: ", logger.Info)
	// zap.S().Info("ori_loggger : ", ori_loggger)
	// zap.S().Info("gormZaplogger : ", gormZaplogger)
	if gin.Mode() == gin.ReleaseMode {
		config = &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			// SkipDefaultTransaction:                   true,
		}
	} else {
		config = &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			// SkipDefaultTransaction:                   true,
			Logger: gormZaplogger,
		}
	}

	//注意：User和Password为MySQL資料庫的管理員密碼，Host和Port為資料庫連接ip端口，DBname為要連接的資料庫
	// 使用配置對象的 DSN 方法生成連接字符串
	mysqlConfig := serverEnv.MysqlDB
	dsn := mysqlConfig.DSN()
	fmt.Printf("正在連接 MySQL: %s\n", mysqlConfig.String())

	var err error
	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return fmt.Errorf("mysql connect failed, err: %w", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(constset.MySQLMaxIdleConns)
	sqlDB.SetMaxOpenConns(constset.MySQLMaxOpenConns)
	m.db = db
	fmt.Println("mysql connect success")
	return nil
}

// Close 關閉 MySQL 連接
func (m *MysqlDB) Close() error {
	// 1. 獲取底層的通用數據庫對象 sql.DB
	sqlDB, err := m.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get mysql db: %w", err)
	}

	// 2. 呼叫 Close()
	return sqlDB.Close()
}

// GetDB 返回底層的 *gorm.DB 對象
func (m *MysqlDB) GetDB() interface{} {
	return m.db
}
