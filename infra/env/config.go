// Package env 負責應用程式的環境配置管理，提供從 YAML 文件載入配置、驗證配置完整性以及線程安全的配置訪問和更新功能。
package env

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

// mapstructure 是用来讀取 yaml 文件字段名 tag

// MysqlConfig Mysql 數據庫配置
type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	DBName   string `mapstructure:"dbName" json:"dbName"`
	Username string `mapstructure:"userName" json:"userName"`
	Password string `mapstructure:"password" json:"password"`
}

// Validate 驗證 MySQL 配置
func (c *MysqlConfig) Validate() error {
	if c.Host == "" {
		return errors.New("mysql host is required")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return errors.New("mysql port must be between 1 and 65535")
	}
	if c.DBName == "" {
		return errors.New("mysql database name is required")
	}
	if c.Username == "" {
		return errors.New("mysql username is required")
	}
	return nil
}

// DSN 生成 MySQL 連接字符串
func (c *MysqlConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.DBName)
}

// String 實現 Stringer 接口（脫敏密碼）
func (c *MysqlConfig) String() string {
	return fmt.Sprintf("MysqlConfig{Host:%s, Port:%d, DBName:%s, Username:%s, Password:***}",
		c.Host, c.Port, c.DBName, c.Username)
}

// MongoDBConfig  MongoDB 數據庫配置
type MongoDBConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	DBName   string `mapstructure:"dbName" json:"dbName"`
	Username string `mapstructure:"userName" json:"userName"`
	Password string `mapstructure:"password" json:"password"`
}

// Validate 驗證 MongoDB 配置
func (c *MongoDBConfig) Validate() error {
	if c.Host == "" {
		return errors.New("mongodb host is required")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return errors.New("mongodb port must be between 1 and 65535")
	}
	return nil
}

// String 實現 Stringer 接口（脫敏密碼）
func (c *MongoDBConfig) String() string {
	return fmt.Sprintf("MongoDBConfig{Host:%s, Port:%d, DBName:%s, Username:%s, Password:***}",
		c.Host, c.Port, c.DBName, c.Username)
}

// RedisConfig  Redis 數據庫配置
type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	DBName   string `mapstructure:"dbName" json:"dbName"`
	Password string `mapstructure:"password" json:"password"`
}

// Validate 驗證 Redis 配置
func (c *RedisConfig) Validate() error {
	if c.Host == "" {
		return errors.New("redis host is required")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return errors.New("redis port must be between 1 and 65535")
	}
	return nil
}

// String 實現 Stringer 接口（脫敏密碼）
func (c *RedisConfig) String() string {
	return fmt.Sprintf("RedisConfig{Host:%s, Port:%d, DBName:%s, Password:***}",
		c.Host, c.Port, c.DBName)
}

// ServerConfig 服務器配置
type ServerConfig struct {
	AppMode   string        `mapstructure:"APP_Mode" json:"APP_Mode"`
	Port      int           `mapstructure:"Port" json:"Port"`
	JwtSecret string        `mapstructure:"JwtSecret" json:"JwtSecret"`
	MysqlDB   MysqlConfig   `mapstructure:"Mysql" json:"Mysql"`
	Redis     RedisConfig   `mapstructure:"Redis" json:"Redis"`
	MongoDB   MongoDBConfig `mapstructure:"MongoDB" json:"MongoDB"`
}

// NewServerConfig 創建一個空的服務器配置實例
// 注意：配置值應從 env.yaml 文件中載入，此函數僅創建空結構
func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}

// Validate 驗證服務器配置的完整性
func (c *ServerConfig) Validate() error {
	// 驗證 AppMode
	validModes := []string{"debug", "release", "test"}
	isValidMode := false
	for _, mode := range validModes {
		if c.AppMode == mode {
			isValidMode = true
			break
		}
	}
	if !isValidMode {
		return fmt.Errorf("invalid app mode: %s, must be one of: %s", c.AppMode, strings.Join(validModes, ", "))
	}

	// 驗證 Port
	if c.Port <= 0 || c.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}

	// 驗證 JwtSecret
	if c.JwtSecret == "" {
		return errors.New("jwt secret is required")
	}
	if len(c.JwtSecret) < 8 {
		return errors.New("jwt secret must be at least 8 characters")
	}

	// 驗證數據庫配置
	if err := c.MysqlDB.Validate(); err != nil {
		return fmt.Errorf("mysql config error: %w", err)
	}

	// Redis 和 MongoDB 可選，只在使用時驗證
	// if err := c.Redis.Validate(); err != nil {
	// 	return fmt.Errorf("redis config error: %w", err)
	// }

	return nil
}

// String 實現 Stringer 接口（脫敏敏感信息）
func (c *ServerConfig) String() string {
	return fmt.Sprintf(`ServerConfig{
	AppMode: %s
	Port: %d
	JwtSecret: ***
	MySQL: %s
	Redis: %s
	MongoDB: %s
}`, c.AppMode, c.Port, c.MysqlDB.String(), c.Redis.String(), c.MongoDB.String())
}

// ConfigManager 配置管理器（線程安全）
type ConfigManager struct {
	mu     sync.RWMutex
	config *ServerConfig
}

// NewConfigManager 創建新的 ConfigManager 實例
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		mu:     sync.RWMutex{},
		config: NewServerConfig(),
	}
}

// GetServerEnv 安全地獲取配置（併發讀）
func (c *ConfigManager) GetServerEnv() *ServerConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}

// UpdateConfig 更新配置（併發寫）
func (c *ConfigManager) UpdateConfig(newConfig *ServerConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config = newConfig
	fmt.Println("配置已更新:", newConfig)
}
