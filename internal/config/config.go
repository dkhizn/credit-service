// Package config содержит структуру конфигурации приложения.
package config

// Config описывает общую конфигурацию приложения.
type Config struct {
	App  App  `yaml:"app"`
	HTTP HTTP `yaml:"http"`
}

// App содержит настройки приложения.
type App struct {
	Name    string `yaml:"name" env-required:"true" env:"APP_NAME"`
	Version string `yaml:"version" env-required:"true" env:"APP_VERSION"`
}

// HTTP содержит настройки HTTP-сервера.
type HTTP struct {
	Port string `yaml:"port" env-required:"true" env:"HTTP_PORT"`
}
