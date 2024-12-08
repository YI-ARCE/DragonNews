package user

import (
	"yiarce/core/dhttp"
)

// Info 手机号登录
//
//	-param 请求参数
//	<br>&emsp; type int 登录类型,1->密码,2->验证码
//	<br>&emsp; phone string 手机号
//	<br>&emsp; code string 验证码 *
//	<br>&emsp; password string 密码 *
//	-method post
func Info(dn *dhttp.Dn) {
}
