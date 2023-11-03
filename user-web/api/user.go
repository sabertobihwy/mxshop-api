package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop-api/user-web/global/response"
	"mxshop-api/user-web/proto"
	"net/http"
	"time"
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
	zap.S().Debugf("get the user list...")
	PORT := 50051
	IP := "127.0.0.1"
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", IP, PORT), grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		zap.S().Errorw("connect to port error...", "msg", err.Error())
	}
	client := proto.NewUserClient(conn)
	lst, err := client.GetUserList(c, &proto.PageInfo{Pn: 0, PSize: 3})
	if err != nil {
		zap.S().Errorw("invoking [GetUserList] error")
		GrpcCodeToHttp(err, c)
		return
	}
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
