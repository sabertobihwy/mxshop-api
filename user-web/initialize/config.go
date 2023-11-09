package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop-api/user-web/utils"

	"mxshop-api/user-web/global"
)

func GetSystemEnv(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	flg := GetSystemEnv("MXSHOP_CONFIG_FLAG")
	filePrefix := "config"
	configFilePath := fmt.Sprintf("user-web/%s_pro.yaml", filePrefix)
	if flg {
		configFilePath = fmt.Sprintf("user-web/%s_local.yaml", filePrefix)
	}
	zap.S().Infof(configFilePath)

	v := viper.New()
	v.SetConfigFile(configFilePath)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(global.SrvConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("SrvConfig : %v", *global.SrvConfig)

	if !flg { // local: fixed port; remote: dynamic port
		global.SrvConfig.Port, _ = utils.GetFreePort()
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("config file changed")
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.SrvConfig)
		zap.S().Infof("SrvConfig change to : %v", *global.SrvConfig)
	})

}
