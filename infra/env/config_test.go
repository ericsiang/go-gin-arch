package env

import (
	"strings"
	"testing"
)

func TestNewServerConfig(t *testing.T) {
	config := NewServerConfig()

	// 驗證新創建的配置是空的（不包含預設值）
	if config.Port != 0 {
		t.Errorf("Expected Port to be 0 (empty), got %d", config.Port)
	}

	if config.AppMode != "" {
		t.Errorf("Expected AppMode to be empty, got %s", config.AppMode)
	}

	if config.JwtSecret != "" {
		t.Errorf("Expected JwtSecret to be empty, got %s", config.JwtSecret)
	}

	t.Log("NewServerConfig correctly creates empty config instance")
}

func TestServerConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *ServerConfig
		wantErr bool
	}{
		{
			name: "有效配置",
			config: &ServerConfig{
				AppMode:   "debug",
				Port:      8080,
				JwtSecret: "test_secret_key_12345",
				MysqlDB: MysqlConfig{
					Host:     "localhost",
					Port:     3306,
					DBName:   "testdb",
					Username: "root",
					Password: "password",
				},
			},
			wantErr: false,
		},
		{
			name: "無效的 AppMode",
			config: &ServerConfig{
				AppMode:   "invalid_mode",
				Port:      8080,
				JwtSecret: "test_secret_key_12345",
				MysqlDB: MysqlConfig{
					Host:     "localhost",
					Port:     3306,
					DBName:   "testdb",
					Username: "root",
					Password: "password",
				},
			},
			wantErr: true,
		},
		{
			name: "無效的 Port",
			config: &ServerConfig{
				AppMode:   "debug",
				Port:      99999,
				JwtSecret: "test_secret_key_12345",
				MysqlDB: MysqlConfig{
					Host:     "localhost",
					Port:     3306,
					DBName:   "testdb",
					Username: "root",
					Password: "password",
				},
			},
			wantErr: true,
		},
		{
			name: "JWT Secret 太短",
			config: &ServerConfig{
				AppMode:   "debug",
				Port:      8080,
				JwtSecret: "short",
				MysqlDB: MysqlConfig{
					Host:     "localhost",
					Port:     3306,
					DBName:   "testdb",
					Username: "root",
					Password: "password",
				},
			},
			wantErr: true,
		},
		{
			name: "MySQL 配置缺少 Host",
			config: &ServerConfig{
				AppMode:   "debug",
				Port:      8080,
				JwtSecret: "test_secret_key_12345",
				MysqlDB: MysqlConfig{
					Host:     "",
					Port:     3306,
					DBName:   "testdb",
					Username: "root",
					Password: "password",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ServerConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  MysqlConfig
		wantErr bool
	}{
		{
			name: "有效的 MySQL 配置",
			config: MysqlConfig{
				Host:     "localhost",
				Port:     3306,
				DBName:   "testdb",
				Username: "root",
				Password: "password",
			},
			wantErr: false,
		},
		{
			name: "缺少 Host",
			config: MysqlConfig{
				Host:     "",
				Port:     3306,
				DBName:   "testdb",
				Username: "root",
				Password: "password",
			},
			wantErr: true,
		},
		{
			name: "無效的 Port",
			config: MysqlConfig{
				Host:     "localhost",
				Port:     0,
				DBName:   "testdb",
				Username: "root",
				Password: "password",
			},
			wantErr: true,
		},
		{
			name: "缺少 DBName",
			config: MysqlConfig{
				Host:     "localhost",
				Port:     3306,
				DBName:   "",
				Username: "root",
				Password: "password",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("MysqlConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlConfigDSN(t *testing.T) {
	config := MysqlConfig{
		Host:     "localhost",
		Port:     3306,
		DBName:   "testdb",
		Username: "root",
		Password: "password123",
	}

	expected := "root:password123@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	actual := config.DSN()

	if actual != expected {
		t.Errorf("MysqlConfig.DSN() = %v, want %v", actual, expected)
	}
}

func TestConfigString(t *testing.T) {
	config := &ServerConfig{
		AppMode:   "debug",
		Port:      8080,
		JwtSecret: "secret_password_12345",
		MysqlDB: MysqlConfig{
			Host:     "localhost",
			Port:     3306,
			DBName:   "testdb",
			Username: "root",
			Password: "password",
		},
	}

	str := config.String()

	// 確保密碼被脫敏
	if strings.Contains(str, "password") {
		t.Error("ServerConfig.String() should mask password")
	}
	if strings.Contains(str, "secret_password") {
		t.Error("ServerConfig.String() should mask JWT secret")
	}

	// 確保包含基本信息
	if !strings.Contains(str, "debug") {
		t.Error("ServerConfig.String() should contain AppMode")
	}
	if !strings.Contains(str, "8080") {
		t.Error("ServerConfig.String() should contain Port")
	}
}

func TestRedisConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  RedisConfig
		wantErr bool
	}{
		{
			name: "有效的 Redis 配置",
			config: RedisConfig{
				Host:     "localhost",
				Port:     6379,
				DBName:   "0",
				Password: "",
			},
			wantErr: false,
		},
		{
			name: "缺少 Host",
			config: RedisConfig{
				Host:     "",
				Port:     6379,
				DBName:   "0",
				Password: "",
			},
			wantErr: true,
		},
		{
			name: "無效的 Port",
			config: RedisConfig{
				Host:     "localhost",
				Port:     99999,
				DBName:   "0",
				Password: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
