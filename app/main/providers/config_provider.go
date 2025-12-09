package providers

import (
	"github.com/brunobotter/notification-system/app/main/config"
	"github.com/brunobotter/notification-system/app/main/container"
	"github.com/brunobotter/notification-system/app/main/logger"
)

type ConfigServiceProvider struct{}

func NewConfigServiceProvider() *ConfigServiceProvider {
	return &ConfigServiceProvider{}
}

func (p *ConfigServiceProvider) Register(c container.Container) {
	c.Singleton(func() *config.Config {
		cfg := config.Init()
		return cfg
	})

	c.Singleton(func(cfg *config.Config) logger.Logger {
		return logger.NewLoggerZap(cfg.App_Name)
	})
}
