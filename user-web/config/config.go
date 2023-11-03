package config

type UserConfig struct {
	Port int32 `mapstructure:"port"`
}
type ServiceConfig struct {
	Ip         string     `mapstructure:"ip"`
	Name       string     `mapstructure:"name"`
	Port       int32      `mapstructure:"port"`
	UserConfig UserConfig `mapstructure:"usr_srv"`
}
