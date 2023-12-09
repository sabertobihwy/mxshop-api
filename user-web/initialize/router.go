package initialize

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/user-web/middlewares"
)
import "mxshop-api/user-web/router"

func InitRouter() *gin.Engine {
	zap.S().Infof("init the router...")
	router := gin.Default()
	router.Use(middlewares.Cors())
	group := router.Group("/v1")
	router2.UserRouter(group)
	router2.BaseRouter(group)
	return router
}
