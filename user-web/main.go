package main

import (
	"fmt"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
)

func main() {
	// 1.init logger
	initialize.InitLogger()
	// 2. init config
	initialize.InitConfig()
	port := global.SrvConfig.Port
	// 3. init router
	router := initialize.InitRouter()
	zap.S().Debugf("start the server...port:%d", port)
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panicf("starting failed...port:%d, error:%s", port, err.Error())

	}
}
