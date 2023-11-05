package router2

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/user-web/api"
)

func InitBaseRouter(group *gin.RouterGroup) {
	zap.S().Infof("init the BaseRouter...")
	group = group.Group("/base")
	{
		group.GET("/captcha", api.GetCaptcha)
	}
}
