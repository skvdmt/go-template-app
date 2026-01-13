package model

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	// Название приложения.
	APP_NAME = "go-template-app"

	// Путь в директории конфигурации. (Добавляется директория с именем приложения).
	configDirectory = "/etc"
	// Имя файла конфигурации.
	configFileName = "config.yaml"
)

// Config Глобальная конфигурация.
var Config *MainConfig

// Postgres Конфигурация соединения с postgres.
type Postgres struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	User     string `yaml:"user"`
	Database string `yaml:"database"`
}

// MainConfig Основная конфигурация.
type MainConfig struct {
	Postgres *Postgres `yaml:"postgres"`
}

// LoadConfig Загрузка конфигурации в глобальную переменную Config.
func LoadConfig() error {
	d, err := os.ReadFile(filepath.Join(configDirectory, APP_NAME, configFileName))
	if err != nil {
		return err
	}
	Config := &MainConfig{}
	if err := yaml.Unmarshal(d, Config); err != nil {
		return err
	}
	return nil
}
