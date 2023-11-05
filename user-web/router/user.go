package router2

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/user-web/api"
	"mxshop-api/user-web/middlewares"
)

func UserRouter(group *gin.RouterGroup) {
	zap.S().Infof("init the UserRouter...")
	group = group.Group("/user")
	{
		group.GET("/list", middlewares.JWTAuth(), middlewares.IsAdmin(), api.GetUserList)
		group.POST("/login", api.LoginValidate)
		group.POST("/register", api.Register)
	}
}
