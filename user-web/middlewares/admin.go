package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/user-web/models"
	"net/http"
)

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, _ := c.Get("claims")
		customClaim, _ := claim.(*models.CustomClaims)
		zap.S().Infof("Is admin? %s:%d", customClaim.NickName, customClaim.AuthorityId)
		if customClaim.AuthorityId != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "not admin",
			})
			c.Abort()
		}
		c.Next()
	}
}
