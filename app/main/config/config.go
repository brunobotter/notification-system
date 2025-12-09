package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func Init() *Config {
	cfg, err := Read()
	if err != nil {
		panic(fmt.Sprintf("Erro ao ler configuração de arquivo: %v", err))
	}
	return cfg
}

func Read() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.BindEnv("server.port", "SERVER_PORT")
	v.BindEnv("server.host", "SERVER_HOST")
	v.BindEnv("app_name", "APP_NAME")
	v.BindEnv("env", "env")

	err := v.ReadInConfig()
	if err != nil && errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return nil, err
	}

	conf := Config{}
	err = v.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil

}
