package api

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/rand"
)

func generateSMSCode(len int) string {
	array := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	rand.Seed(uint64(time.Now().UnixNano()))
	var sb strings.Builder
	for i := 0; i < len; i++ {
		sb.WriteString(fmt.Sprintf("%d", array[rand.Intn(10)]))
	}
	return sb.String()
}

func SendSMS(c *gin.Context) {
	// 1, validate the form
	var sms = forms.SendSmsForm{}
	if err := c.ShouldBind(&sms); err != nil {
		validateReturn(err, c)
		return
	}
	// 2, mock receive sms code, save for 15s
	smsCode := generateSMSCode(5)
	zap.S().Infof("generate sms code : %s", smsCode)
	zap.S().Infof("redis: %s", fmt.Sprintf("%s:%d", global.SrvConfig.RedisConfig.Host, global.SrvConfig.RedisConfig.Port))
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.SrvConfig.RedisConfig.Host, global.SrvConfig.RedisConfig.Port),
	})
	rdb.Set(c, sms.Mobile, smsCode, 300*time.Second)

	c.JSON(http.StatusOK, gin.H{
		"msg": "send sms_code successfully",
	})
}
