package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/consul/api"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"

	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
	"mxshop-api/user-web/utils"
	myvalidator "mxshop-api/user-web/validator"
)

var serviceId string

func Register(consuladdr string, httpIp string, port int, name string, tags []string) (*api.Client, error) {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", consuladdr, 8500)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	// health check instance
	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/v1/base/health", httpIp, port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}
	serviceId = fmt.Sprintf("%s", uuid.NewV4())
	// registration instance
	regis := &api.AgentServiceRegistration{
		Name:    name,
		ID:      serviceId,
		Port:    port,
		Tags:    tags,
		Address: httpIp,
		Check:   check,
	}
	err = client.Agent().ServiceRegister(regis)
	if err != nil {
		panic(err)
	}
	return client, err

}

func main() {
	// 1.init logger
	initialize.InitLogger()
	// 2. init config
	initialize.InitConfig()
	//port := global.SrvConfig.Port
	port, _ := utils.GetFreePort()
	// 3. init trans
	if err := initialize.InitTrans("en"); err != nil {
		zap.S().Debugf("init translator errror,%s", err.Error())
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}
	// 3. init client
	initialize.InitUsrSrv()
	defer global.Conn.Close()
	// 4. init router
	router := initialize.InitRouter()
	// 5. init sentinel
	initialize.InitializeSentinel()
	// 6. register consul
	client, err := Register(global.SrvConfig.ConsulConfig.Host, global.SrvConfig.Ip, port,
		global.SrvConfig.Name, []string{"mxshop", "bobby"})
	if err != nil {
		zap.S().Debugf(err.Error())
		panic("register error")
	}
	zap.S().Debugf("start the server...port:%d", port)
	go func() {
		if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
			zap.S().Panicf("starting failed...port:%d, error:%s", port, err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(serviceId); err != nil {
		zap.S().Debugf(fmt.Sprintf("deregister error"))
		panic(err)
	}
}
