package utils

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

func GetService(consulAddr string, consulPort int, srvName string) (srvAddr string, srvPort int) {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", consulAddr, consulPort)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, srvName))
	if err != nil {
		zap.S().Debugf("cannot find the service")
		panic(err)
	}
	return data[srvName].Address, data[srvName].Port
}
