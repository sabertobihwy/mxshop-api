package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/proto"
)

func getService(consulAddr string, consulPort int, srvName string) (srvAddr string, srvPort int) {
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
	if data[srvName].Address == "" {
		zap.S().Fatalf("connect to [%s] error", srvName)
		panic("connect to service error")
	}
	return data[srvName].Address, data[srvName].Port
}

func InitUsrSrv() {
	SRV_HOST, SRV_PORT := getService(global.SrvConfig.ConsulConfig.Host,
		global.SrvConfig.ConsulConfig.Port, "mxshop_srvs")
	zap.S().Infof("SRV_HOST: %s SRV_PORT: %d", SRV_HOST, SRV_PORT)
	var err error
	global.Conn, err = grpc.Dial(fmt.Sprintf("%s:%d", SRV_HOST, SRV_PORT), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("connect to port error...", "msg", err.Error())
	}
	global.UserClient = proto.NewUserClient(global.Conn)
}
