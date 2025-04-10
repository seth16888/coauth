package model

// PasswordLoginReq 密码登录请求参数
type PasswordLoginReq struct {
	Username    string `json:"username" form:"username" validate:"required,min=4,max=32"`
	Password    string `json:"password" form:"password" validate:"required,min=6,max=32"`
	CaptchaKey  string `json:"captchaKey" form:"captchaKey" validate:"required,min=12,max=64"`
	CaptchaCode string `json:"captchaCode" form:"captchaCode" validate:"required,min=4,max=8"`
}

// PasswordLoginRes 登录返回
type PasswordLoginRes struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int64  `json:"expiresIn"`
}

// VerifyCaptchaReq 验证验证码请求参数
type VerifyCaptchaReq struct {
	CaptchaKey  string `json:"captchaKey" validate:"required,min=12,max=64"`
	CaptchaCode string `json:"captchaCode" validate:"required,min=4,max=8"`
	Clear       bool   `json:"clear"` // 是否清除验证码
}
