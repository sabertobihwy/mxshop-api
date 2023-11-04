package global

import (
	ut "github.com/go-playground/universal-translator"
	"mxshop-api/user-web/config"
)

var (
	SrvConfig *config.ServiceConfig = &config.ServiceConfig{}
	Trans     ut.Translator
)
