package initialize

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)
import "mxshop-api/user-web/router"

func InitRouter() *gin.Engine {
	zap.S().Infof("init the router...")
	router := gin.Default()
	group := router.Group("/u/v1")
	router2.UserRouter(group)
	return router
}
func InitLogger() {
	// set a global logger to use zap.S() -> secure access by goroutines with Mutex lock
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}
