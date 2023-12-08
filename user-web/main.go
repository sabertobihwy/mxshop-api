package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
	myvalidator "mxshop-api/user-web/validator"
)

func main() {
	// 1.init logger
	initialize.InitLogger()
	// 2. init config
	initialize.InitConfig()
	port := global.SrvConfig.Port
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
	zap.S().Debugf("start the server...port:%d", port)
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panicf("starting failed...port:%d, error:%s", port, err.Error())

	}
}
