package index

import (
	"yiarce/core/dhttp"
	"yiarce/core/frame"
)

func Hello(d *dhttp.Dn) {
	d.Write(200, `Hello DragonNews`)
}

func GetKey(d *dhttp.Dn) {
	d.Json(d.GetAll())
}

func GetJson(d *dhttp.Dn) {
	// 如果提交数据复杂,请使用Format
	//data := map[string]interface{}{}
	//d.BodyFormat(data, `服务器错误`)
	// ------------------------
	// 返回map[string]string类型
	frame.Println(d.Data())
	d.Json(d.Data())
}
