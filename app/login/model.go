package login

// PhoneModel 手机号登录
type PhoneModel struct {
	Type     int    `json:"type"`               // 登录类型,1->密码,2->验证码
	Phone    string `json:"phone"`              // 手机号
	Code     string `json:"code,omitempty"`     // 验证码
	Password string `json:"password,omitempty"` // 密码
}

// Phone2Model 手机号登录2
type Phone2Model struct {
	Type int `json:"type"` // 登录类型,1->密码,2->验证码
}
