package global

import (
	ut "github.com/go-playground/universal-translator"
	"mxshop-api/user-web/config"
	"mxshop-api/user-web/proto"

	"google.golang.org/grpc"
)

var (
	SrvConfig  *config.ServiceConfig = &config.ServiceConfig{}
	Trans      ut.Translator
	Conn       *grpc.ClientConn
	UserClient proto.UserClient
)
