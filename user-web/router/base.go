package router2

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/user-web/api"
)

func BaseRouter(group *gin.RouterGroup) {
	zap.S().Infof("init the BaseRouter...")
	group = group.Group("/base")
	{
		group.GET("/captcha", api.GetCaptcha)
		group.POST("/send_sms", api.SendSMS)
	}
}
