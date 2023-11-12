package config

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}
type JwtKey struct {
	Key string `mapstructure:"key" json:"key"`
}
type RedisConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}
type UserServiceConfig struct {
	Name string `mapstructure:"name" json:"name"`
}
type ServiceConfig struct {
	Ip                string            `mapstructure:"ip" json:"ip"`
	Name              string            `mapstructure:"name" json:"name"`
	Port              int               `mapstructure:"port" json:"port"`
	ConsulConfig      ConsulConfig      `mapstructure:"consul" json:"consul"`
	JwtConfig         JwtKey            `mapstructure:"jwt" json:"jwt"`
	RedisConfig       RedisConfig       `mapstructure:"redis" json:"redis"`
	UserServiceConfig UserServiceConfig `mapstructure:"user_srv" json:"user_srv"`
}
type NacosConfig struct {
	Host      string `mapstructure:"host" json:"host"`
	Port      int    `mapstructure:"port" json:"port"`
	Namespace string `mapstructure:"namespace" json:"namespace"`
	Dataid    string `mapstructure:"dataid" json:"dataid"`
	Group     string `mapstructure:"group" json:"group"`
	User      string `mapstructure:"user" json:"user"`
	Password  string `mapstructure:"password" json:"password"`
}
