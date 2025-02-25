package core

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Config 结构体，用于存储配置信息
type Config struct {
	DefaultSSH SSHConfig   `json:"default_ssh"`  // 默认SSH连接配置
	Servers    []SSHConfig `json:"servers"`      // 多服务器配置列表
	Auth       AuthConfig  `json:"auth"`         // 软件认证配置
}

// SSHConfig 结构体，用于存储SSH连接配置
type SSHConfig struct {
	Name         string `json:"name"`           // 服务器名称
	Host         string `json:"host"`           // 主机地址
	Port         int    `json:"port"`           // 端口
	Username     string `json:"username"`       // 用户名
	Password     string `json:"password"`       // 密码
	LoginType    int    `json:"login_type"`     // 登录类型（0: 密码认证, 1: 私钥认证）
	KeyPassphrase string `json:"key_passphrase"` // 私钥密码
}

// AuthConfig 结构体，用于存储软件认证配置
type AuthConfig struct {
	Enable   bool   `json:"enable"`    // 是否启用认证
	Username string `json:"username"`  // 认证用户名
	Password string `json:"password"`  // 认证密码
}

// 默认配置
var defaultConfig = Config{
	DefaultSSH: SSHConfig{
		Name:          "Default",
		Host:          "",
		Port:          22,
		Username:      "root",
		Password:      "",
		LoginType:     0,
		KeyPassphrase: "",
	},
	Servers: []SSHConfig{},
	Auth: AuthConfig{
		Enable:   false,
		Username: "",
		Password: "",
	},
}

// LoadConfig 加载配置文件
func LoadConfig() Config {
	// 获取可执行文件所在目录
	exePath, err := os.Executable()
	if err != nil {
		log.Println("Failed to get executable path:", err)
		return defaultConfig
	}
	
	exeDir := filepath.Dir(exePath)
	configPath := filepath.Join(exeDir, "config.json")
	
	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Println("Config file not found, using default settings")
		return defaultConfig
	}
	
	// 读取配置文件
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Println("Failed to read config file:", err)
		return defaultConfig
	}
	
	// 解析配置文件
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Println("Failed to parse config file:", err)
		return defaultConfig
	}
	
	// 设置默认值（如果配置文件中没有指定）
	if config.DefaultSSH.Port == 0 {
		config.DefaultSSH.Port = 22
	}
	
	if config.DefaultSSH.Username == "" {
		config.DefaultSSH.Username = "root"
	}
	
	// 确保每个服务器配置都有合理的默认值
	for i := range config.Servers {
		if config.Servers[i].Port == 0 {
			config.Servers[i].Port = 22
		}
		
		if config.Servers[i].Username == "" {
			config.Servers[i].Username = "root"
		}
		
		// 如果没有设置名称，使用主机地址作为名称
		if config.Servers[i].Name == "" {
			config.Servers[i].Name = config.Servers[i].Host
		}
	}
	
	log.Println("Config loaded successfully")
	return config
}

// GetDefaultSSHConfig 获取默认SSH配置
func GetDefaultSSHConfig() SSHConfig {
	config := LoadConfig()
	return config.DefaultSSH
}

// GetServerList 获取服务器列表
func GetServerList() []SSHConfig {
	config := LoadConfig()
	return config.Servers
}

// GetAuthConfig 获取认证配置
func GetAuthConfig() AuthConfig {
	config := LoadConfig()
	return config.Auth
} 