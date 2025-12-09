package config

type Config struct {
	Server   ServerConfig `mapstruct:"server"`
	App_Name string       `mapstruct:"app_name"`
	Env      string       `mapstruct:"env"`
}

type ServerConfig struct {
	Port int    `mapstruct:"port"`
	Host string `mapstruct:"host"`
}
