package forms

type Login struct {
	Mobile    string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Password  string `form:"password" json:"password" binding:"required,min=3,max=10"`
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"`
	CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required"`
}

type RegisterForm struct {
	NickName string `form:"nickName" json:"nickName" binding:"required"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=10"`
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"`
	SmsCode  string `form:"smscode" json:"smscode" binding:"required"`
}
