package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"net/http"
)

var store = base64Captcha.DefaultMemStore

func GetCaptcha(c *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := cp.Generate()
	if err != nil {
		zap.S().Errorf("generate captcha error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "generate captcha error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"captcha_id":   id,
		"captcha_path": b64s,
	})
}
