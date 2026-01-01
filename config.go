package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	ConsulAddr  string `yaml:"consul_addr"`
	ServiceName string `yaml:"service_name"`
	HTTPPort    string `yaml:"http_port"`
	ServiceAddr string `yaml:"service_addr"` // 要调用的微服务名称，默认 "edu-course"
}

// ConfigFile YAML 配置文件结构
type ConfigFile struct {
	Consul struct {
		Addr string `yaml:"addr"`
	} `yaml:"consul"`
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Service struct {
		Addr string `yaml:"addr"`
		Name string `yaml:"name"`
	} `yaml:"service"`
}

// LoadConfig 从配置文件、环境变量或默认值加载配置
func LoadConfig() *Config {
	config := &Config{}

	// 尝试从配置文件加载
	configFile := loadConfigFile()
	if configFile != nil {
		config.ConsulAddr = configFile.Consul.Addr
		config.HTTPPort = configFile.Server.Port
		config.ServiceAddr = configFile.Service.Addr
		config.ServiceName = configFile.Service.Name
	}

	// 设置默认值
	if config.ConsulAddr == "" {
		config.ConsulAddr = "http://47.113.220.106:8500"
	}
	if config.HTTPPort == "" {
		config.HTTPPort = "8080"
	}
	if config.ServiceAddr == "" {
		config.ServiceAddr = "edu-course"
	}

	// 环境变量可以覆盖配置文件（可选）
	if envConsul := os.Getenv("CONSUL_ADDR"); envConsul != "" {
		config.ConsulAddr = envConsul
	}
	if envPort := os.Getenv("HTTP_PORT"); envPort != "" {
		config.HTTPPort = envPort
	}
	if envServiceAddr := os.Getenv("SERVICE_ADDR"); envServiceAddr != "" {
		config.ServiceAddr = envServiceAddr
	}
	if envServiceName := os.Getenv("SERVICE_NAME"); envServiceName != "" {
		config.ServiceName = envServiceName
	}

	return config
}

// loadConfigFile 从 config.yaml 文件加载配置
func loadConfigFile() *ConfigFile {
	configPath := "config.yaml"

	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 配置文件不存在时返回 nil，使用默认值
		return nil
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Warning: Failed to read config file %s: %v, using defaults\n", configPath, err)
		return nil
	}

	// 解析 YAML
	var configFile ConfigFile
	if err := yaml.Unmarshal(data, &configFile); err != nil {
		fmt.Printf("Warning: Failed to parse config file %s: %v, using defaults\n", configPath, err)
		return nil
	}

	return &configFile
}

// getEnv 获取环境变量，如果不存在则返回默认值（保留以兼容旧代码）
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
