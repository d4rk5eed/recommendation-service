package models

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type ServiceConfig map[string]map[string]string

func ReadConfig(confPath string) (ServiceConfig, error) {
	confPath, err := filepath.Abs(confPath)
	if err != nil {
		return nil, err
	}

	configBytes, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}

	conf := new(Config)

	// configStr := expandEnv(string(configBytes))
	configStr := string(configBytes)
	if err := yaml.Unmarshal([]byte(configStr), conf); err != nil {
		return nil, err
	}

	// Создаем результирующий словарь
	result := make(map[string]map[string]string)
	for _, endpoint := range conf.Endpoints {
		result[endpoint.Name] = map[string]string{
			"name":  endpoint.Name,
			"url":   endpoint.URL,
			"key":   GetEnvWithDefault("API_KEY", endpoint.Key),
			"model": endpoint.Model,
		}
	}

	return result, nil
}

// Вспомогательные функции
func GetEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

type Config struct {
	Endpoints []Endpoint `yaml:"endpoints"`
}

type Endpoint struct {
	Name  string `yaml:"name"`
	URL   string `yaml:"url"`
	Key   string `yaml:"key"`
	Model string `yaml: "model"`
}
