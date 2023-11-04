package config

type UserConfig struct {
	Port int32 `mapstructure:"port"`
}
type JwtKey struct {
	Key string `mapstructure:"key"`
}
type ServiceConfig struct {
	Ip         string     `mapstructure:"ip"`
	Name       string     `mapstructure:"name"`
	Port       int32      `mapstructure:"port"`
	UserConfig UserConfig `mapstructure:"usr_srv"`
	JwtConfig  JwtKey     `mapstructure:"jwt"`
}
