package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Micro MicroConfig
	HTTP  HTTPConfig
}

type MicroConfig struct {
	Registry          string `yaml:"registry"`
	RegistryAddr      string `yaml:"registry_addr"`
	CourseServiceName string `yaml:"course_service_name"`
}

type HTTPConfig struct {
	Address string `yaml:"address"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	// 优先从 config.yaml 加载配置
	configFile := getEnv("CONFIG_FILE", "config.yaml")
	if _, err := os.Stat(configFile); err == nil {
		data, err := os.ReadFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// 如果配置文件中没有值，使用环境变量或默认值
	if cfg.Micro.Registry == "" {
		cfg.Micro.Registry = getEnv("MICRO_REGISTRY", "mdns")
	}
	if cfg.Micro.RegistryAddr == "" {
		cfg.Micro.RegistryAddr = getEnv("MICRO_REGISTRY_ADDRESS", "")
	}
	if cfg.Micro.CourseServiceName == "" {
		cfg.Micro.CourseServiceName = getEnv("COURSE_SERVICE_NAME", "edu-course")
	}
	if cfg.HTTP.Address == "" {
		cfg.HTTP.Address = getEnv("HTTP_ADDRESS", ":8083")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
