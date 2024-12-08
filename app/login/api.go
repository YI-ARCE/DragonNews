package login

import (
	"dragonNews/table/sakuraPost"
	"yiarce/core/dhttp"
)

// Phone 手机号登录
//
//	-param PhoneModel
//	-method post
func Phone(d *dhttp.Dn) {
	result := d.Table(`xxx`).Where(`c = 1`).Find()
	if result.Err() != nil {
		d.Json(result.Err().Error(), 500)
	}
	d.Json(result.Result(), 500)
}

// Phone2 手机号登录2
//
//	-param PhoneModel2
//	-method post
func Phone2(d *dhttp.Dn) {
	result := d.Table(`xxx`).Where(`a = 2`).Find()
	if result.Err() != nil {
		d.Json(result.Err().Error(), 500)
	}
	m := sakuraPost.Structure{}
	result.Format(&m)
	d.Json(m)
}
