package router2

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/user-web/api"
	"mxshop-api/user-web/middlewares"
)

func UserRouter(group *gin.RouterGroup) {
	zap.S().Infof("init the UserRouter...")
	rg := group.Group("/user").Use(middlewares.Tracing())
	{
		rg.GET("/list", middlewares.JWTAuth(), middlewares.IsAdmin(), api.GetUserList)
		rg.POST("/login", api.LoginValidate)
		rg.POST("/register", api.Register)
	}
}
