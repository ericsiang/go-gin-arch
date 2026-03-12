package zaplog

import (
	"log"
	"self_go_gin/util/zap_logger"
	"time"

	"go.uber.org/zap"
)

const (
	// ZapLoggerMaxSize log 大保存容量，單位為Byte
	ZapLoggerMaxSize      = 1 * 1024 * 1024 
	// ZapLoggerMaxCounts log 最大保存個數，單位為個數
	ZapLoggerMaxCounts    = 2  
	// ZapLoggerMaxAge log 最大保存時間，單位為時間單位，不使用時請設定為 -1            
	ZapLoggerMaxAge       = -1   
	// ZapLoggerRotationTime log 切割頻率，單位為時間單位           
	ZapLoggerRotationTime = 24 * time.Hour  
)

// GetZapLogger 獲取 zap Logger 實例
/*
logPath: 日誌文件保存路徑
*/
func GetZapLogger(logPath string) *zap.Logger {
	rotatelogsConfig := &zap_logger.RotatelogsConfig{
		InfoLogPath:   logPath + "info/info_%Y-%m-%d.log",
		ErrorLogPath:  logPath + "error/error_%Y-%m-%d.log",
		MaxSize:       ZapLoggerMaxSize,
		RotationCount: ZapLoggerMaxCounts,
		MaxAge:        ZapLoggerMaxAge,
		RotationTime:  ZapLoggerRotationTime,
	}

	zapLogger, err := zap_logger.NewLogger(rotatelogsConfig)
	if err != nil {
		log.Fatalln("[logger] GetZapLogger() err : ", err)
	}

	return zapLogger.GetZapLogger()
}
