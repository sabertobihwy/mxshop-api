package config

type ConsulConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
type JwtKey struct {
	Key string `mapstructure:"key"`
}
type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
type UserServiceConfig struct {
	Name string `mapstructure:"name"`
}
type ServiceConfig struct {
	Ip                string            `mapstructure:"ip"`
	Name              string            `mapstructure:"name"`
	Port              int               `mapstructure:"port"`
	ConsulConfig      ConsulConfig      `mapstructure:"consul"`
	JwtConfig         JwtKey            `mapstructure:"jwt"`
	RedisConfig       RedisConfig       `mapstructure:"redis"`
	UserServiceConfig UserServiceConfig `mapstructure:"user_srv"`
}
