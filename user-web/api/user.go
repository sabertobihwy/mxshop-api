package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetUserList(c *gin.Context) {
	zap.S().Debugf("get the user list...")
}
