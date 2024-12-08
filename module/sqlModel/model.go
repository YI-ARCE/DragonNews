package sqlModel

//type table struct {
//	name string
//	keys []tableKey
//}

type tableKey struct {
	// 键名
	name string
	// 是否为主键
	pri bool
	// 键类型
	types string
	// 是否可以为null
	isNull bool
	// 备注
	remark string
}
