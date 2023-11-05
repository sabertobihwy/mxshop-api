package middlewares

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/user-web/models"
	"net/http"
)

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, _ := c.Get("claims")
		customClaim, _ := claim.(*models.CustomClaims)
		if customClaim.AuthorityId != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "not admin",
			})
			c.Abort()
		}
		c.Next()
	}
}
