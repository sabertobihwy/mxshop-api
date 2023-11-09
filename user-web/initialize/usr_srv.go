package initialize

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"mxshop-api/user-web/global"
	"mxshop-api/user-web/proto"
)

func manualWayGetCli() {
	var err error
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.SrvConfig.ConsulConfig.Host, global.SrvConfig.ConsulConfig.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.SrvConfig.UserServiceConfig.Name))
	if err != nil {
		zap.S().Debugf("cannot find the service")
		panic(err)
	}
	if data[global.SrvConfig.UserServiceConfig.Name].Address == "" {
		zap.S().Fatalf("connect to [%s] error", global.SrvConfig.UserServiceConfig.Name)
		panic("connect to service error")
	}
	SRV_HOST := data[global.SrvConfig.UserServiceConfig.Name].Address
	SRV_PORT := data[global.SrvConfig.UserServiceConfig.Name].Port

	global.Conn, err = grpc.Dial(fmt.Sprintf("%s:%d", SRV_HOST, SRV_PORT), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("connect to port error...", "msg", err.Error())
	}
	global.UserClient = proto.NewUserClient(global.Conn)
}

func LBwayGetCli() {
	var err error
	global.Conn, err = grpc.Dial(fmt.Sprintf("consul://%s:8500/%s?wait=14s&tag=%s", global.SrvConfig.ConsulConfig.Host, global.SrvConfig.UserServiceConfig.Name, "mxshop"),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	global.UserClient = proto.NewUserClient(global.Conn)
}

func InitUsrSrv() {
	LBwayGetCli()

}
