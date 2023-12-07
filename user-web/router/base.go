package router2

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/user-web/api"
	"mxshop-api/user-web/middlewares"
)

func BaseRouter(group *gin.RouterGroup) {
	zap.S().Infof("init the BaseRouter...")
	rg := group.Group("/base").Use(middlewares.Tracing())
	{
		rg.GET("/captcha", api.GetCaptcha)
		rg.POST("/send_sms", api.SendSMS)
		rg.GET("/health", api.HealthCheck)
	}
}
