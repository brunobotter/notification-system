package providers

import (
	"github.com/brunobotter/notification-system/infra/logger"
	"github.com/brunobotter/notification-system/main/config"
	"github.com/brunobotter/notification-system/main/container"
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
