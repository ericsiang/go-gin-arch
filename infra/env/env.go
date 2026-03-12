package env

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	envFile       = "env.yaml"
	configManager = NewConfigManager()
)

// LoadConfig 載入配置文件（不監聽變更）
func LoadConfig(envPath string, serverConfig *ServerConfig) error {
	v := viper.New()
	v.SetConfigFile(envPath + envFile)

	// 讀取配置信息
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("配置文件讀取錯誤: %w", err)
	}

	if err := v.Unmarshal(serverConfig); err != nil {
		return fmt.Errorf("配置文件解析錯誤: %w", err)
	}

	// 驗證配置
	if err := serverConfig.Validate(); err != nil {
		return fmt.Errorf("配置驗證失敗: %w", err)
	}

	return nil
}

// InitEnv 初始化環境變量（帶配置文件監聽）
/*
envPath: 環境變量文件所在路徑
serverConfig: 服務器配置實例指針
reloadFunc: 配置文件改變時的回調函數（接收新配置）
*/
func InitEnv(envPath string) error {
	v := viper.New()
	fmt.Printf("讀取配置文件: %s\n", envPath+envFile)
	if envPath == "" {
		// 如果未通过环境变量指定，默认使用可执行文件目录下的 conf 文件夹
		v.SetConfigName("env")    // 檔名
		v.SetConfigType("yaml")   // 格式
		v.AddConfigPath("./conf") // 執行路徑下的 conf
		v.AddConfigPath("../../conf")
	} else {
		v.SetConfigFile(envPath + envFile)
	}

	// 初始載入配置
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("配置文件讀取錯誤: %w", err)
	}

	if err := v.Unmarshal(configManager.GetServerEnv()); err != nil {
		return fmt.Errorf("配置文件解析錯誤: %w", err)
	}

	// 驗證初始配置
	if err := configManager.GetServerEnv().Validate(); err != nil {
		return fmt.Errorf("配置驗證失敗: %w", err)
	}

	fmt.Println("配置載入成功")

	// 監聽配置文件變更
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("檢測到配置文件變更:", e.Name)

		// 創建新配置實例
		newConfig := configManager.GetServerEnv()
		if err := v.Unmarshal(newConfig); err != nil {
			fmt.Fprintf(os.Stderr, "配置熱重載失敗(解析錯誤): %v\n", err)
			return
		}

		// 驗證新配置
		if err := newConfig.Validate(); err != nil {
			fmt.Fprintf(os.Stderr, "配置熱重載失敗(驗證錯誤): %v\n", err)
			return
		}

		// 函數更新配置
		configManager.UpdateConfig(newConfig)

		fmt.Println("配置熱重載成功")
	})

	return nil
}

func GetConfigManager() *ConfigManager {
	return configManager
}
