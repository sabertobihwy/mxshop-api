package router2

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/user-web/api"
)

func UserRouter(group *gin.RouterGroup) {
	zap.S().Infof("init the UserRouter...")
	group = group.Group("/user")
	{
		group.GET("/list", api.GetUserList)
	}
}
