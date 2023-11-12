package initialize

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
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
	if err := v.Unmarshal(global.NacosConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("NacosConfig : %v", *global.NacosConfig)

	//create ServerConfig
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(global.NacosConfig.Host, uint64(global.NacosConfig.Port),
			constant.WithContextPath("/nacos")),
	}

	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(global.NacosConfig.Namespace),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("tmp/nacos/log"),
		constant.WithCacheDir("tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)

	// create config client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err.Error())
	}
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.Dataid,
		Group:  global.NacosConfig.Group,
	})
	//	fmt.Printf("GetConfig,config : %s", content)
	err = json.Unmarshal([]byte(content), &global.SrvConfig)
	if err != nil {
		panic(err.Error())
	}

	if !flg { // local: fixed port; remote: dynamic port
		global.SrvConfig.Port, _ = utils.GetFreePort()
	}

}
