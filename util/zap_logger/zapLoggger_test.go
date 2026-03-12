package zap_logger

import (
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	// MaxAge and RotationCount cannot be both set  兩者不能同時設置
	// 官方 github 上有說明建議使用 WithRotationCount，要將 MaxAge 設為 -1 比較保險
	config := &RotatelogsConfig{
		InfoLogPath:   "./info.log",
		ErrorLogPath:  "./error.log",
		MaxSize:       1024,
		RotationCount: 7,
		MaxAge:        -1,
		RotationTime:  24 * time.Hour,
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create new logger TestNewLogger() error : %v", err)
	}

	if logger == nil {
		t.Fatal("Expected logger to be not nil TestNewLogger()")
	}
}
