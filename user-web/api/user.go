package api

import (
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global/response"
	"mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/models"
	"mxshop-api/user-web/proto"
	"net/http"
	"strconv"
	"strings"
	"time"

	"mxshop-api/user-web/global"
)

func GrpcCodeToHttp(err error, ctx *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "inner error",
				})
			case codes.InvalidArgument:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "InvalidArgument",
				})
			case codes.Unavailable:
				ctx.JSON(http.StatusServiceUnavailable, gin.H{
					"msg": "Service Unavailable",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": fmt.Sprintf("other errors, %s", e.Code()),
				})
			}
		}
	}
}

func GetUserList(c *gin.Context) {
	// get user_id from token
	claim, _ := c.Get("claims")
	customClaim, _ := claim.(*models.CustomClaims)
	zap.S().Debugf("user_id is %+v", customClaim.ID)
	pn := c.DefaultQuery("pn", "0")
	pSize := c.DefaultQuery("psize", "5")
	pni, _ := strconv.Atoi(pn)
	pns, _ := strconv.Atoi(pSize)

	// for tracing
	newctx := context.WithValue(context.Background(), "ginContext", c)
	// for rate limiting
	e, b := sentinel.Entry("api-test", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		fmt.Println("refuse")
		c.JSON(http.StatusTooManyRequests, gin.H{
			"msg": "too frequent requests!!!", // map[string]string
		})
		return
	}
	lst, err := global.UserClient.GetUserList(newctx, &proto.PageInfo{Pn: uint32(pni), PSize: uint32(pns)})
	if err != nil {
		zap.S().Errorw("invoking [GetUserList] error")
		GrpcCodeToHttp(err, c)
		return
	}
	e.Exit()
	result := make([]interface{}, 0)
	for _, value := range lst.Data {
		//data := make(map[string]interface{})
		usr := response.UserResponse{
			Id:       value.Id,
			Mobile:   value.Mobile,
			Password: value.PassWord,
			NickName: value.NickName,
			Birthday: response.TimeJson(time.Unix(int64(value.BirthDay), 0)),
			//time.Unix(int64(value.BirthDay), 0).Format("2006-01-01"),
			//response.TimeJson(time.Unix(int64(value.BirthDay), 0)),
			// time.Now().Format(time.DateOnly),
		}
		result = append(result, usr)
	}
	c.JSON(http.StatusOK, result)

}
func improveStruct(m map[string]string) map[string]string {
	rsp := map[string]string{}
	for key, value := range m {
		rsp[key[strings.Index(key, ".")+1:]] = value
	}
	return rsp
}
func validateReturn(err error, c *gin.Context) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		// 非validator.ValidationErrors类型错误直接返回
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// validator.ValidationErrors类型错误则进行翻译
	c.JSON(http.StatusBadRequest, gin.H{
		"msg": improveStruct(errs.Translate(global.Trans)), // map[string]string
	})
	return
}

func LoginValidate(c *gin.Context) {
	// 1, validate the form
	var login = forms.Login{}
	if err := c.ShouldBind(&login); err != nil {
		validateReturn(err, c)
		return
	}
	// 1,1 verify the captcha
	if !store.Verify(login.CaptchaId, login.Captcha, true) {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "captcha not verified",
		})
		return
	}
	// 2, interaction with db
	// 2,1 check mobile
	userInfo, err := global.UserClient.GetUserByMobile(c, &proto.MobileRequest{
		Mobile: login.Mobile,
	})
	if err != nil {
		zap.S().Debugf("[error] GetUserByMobile, %s", err.Error())
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "user not found",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "login error",
				})
			}
			return
		}
	} else {
		// 2,2 get the pwd
		// check the pwd
		if rsp, err := global.UserClient.CheckPwd(c, &proto.PwdCheckInfo{
			PassWord:     login.Password,
			EncryptedPws: userInfo.PassWord,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "login  error",
			})
		} else {
			if rsp.Success {
				// 3. create jwt
				j := middlewares.NewJWT()
				claim := models.CustomClaims{
					ID:          uint(userInfo.Id),
					NickName:    userInfo.NickName,
					AuthorityId: uint(userInfo.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),               // start from now
						ExpiresAt: time.Now().Unix() + 60*60*24*30, // 30 days
						Issuer:    "bobby",
					},
				}
				zap.S().Infof("login a user, roleId is :%d", claim.AuthorityId)
				token, err := j.CreateToken(claim)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "create token error",
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"msg":          "login success",
					"token":        token,
					"id":           userInfo.Id,
					"nickname":     userInfo.NickName,
					"token_expire": (time.Now().Unix() + 60*60*24*30) * 1000,
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "password wrong",
				})
			}
		}

	}
}

func Register(c *gin.Context) {
	var registerForm = forms.RegisterForm{}
	if err := c.ShouldBind(&registerForm); err != nil {
		validateReturn(err, c)
		return
	}
	// 2. verify the smscode
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.SrvConfig.RedisConfig.Host, global.SrvConfig.RedisConfig.Port),
	})
	value, err := rdb.Get(c, registerForm.Mobile).Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "sms code error",
		})
		return
	} else {
		if value != registerForm.SmsCode {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": "sms code error",
			})
			return
		}
	}
	// 3. already exists

	if _, err := global.UserClient.GetUserByMobile(c, &proto.MobileRequest{
		Mobile: registerForm.Mobile,
	}); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "user already exists",
		})
		return
	}
	// 4. create new user
	userInfo, err := global.UserClient.CreateUser(c, &proto.CreateUserInfo{
		NickName: registerForm.NickName,
		PassWord: registerForm.Password,
		Mobile:   registerForm.Mobile,
	})
	if err != nil {
		zap.S().Errorw("invoking [CreateNewUser] error")
		GrpcCodeToHttp(err, c)
		return
	}
	// 5. add token
	j := middlewares.NewJWT()
	claim := models.CustomClaims{
		ID:          uint(userInfo.Id),
		NickName:    userInfo.NickName,
		AuthorityId: uint(userInfo.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               // start from now
			ExpiresAt: time.Now().Unix() + 60*60*24*30, // 30 days
			Issuer:    "bobby",
		},
	}
	token, err := j.CreateToken(claim)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "create token error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":          "register success",
		"token":        token,
		"id":           userInfo.Id,
		"nickname":     userInfo.NickName,
		"token_expire": (time.Now().Unix() + 60*60*24*30) * 1000,
	})
}
