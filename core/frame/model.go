package frame

type ErrorPackage struct {
	// 报错的包名
	PackageName string
	// 报错的包名函数
	FuncName string
	// 报错的物理路径
	Path string
	// 报错的文件名
	FileName string
	// 捕获位置
	Line int
}

type Error struct {
	Message string // 报错提示
	// 方法接收的内容,仅限于框架自身
	RequestArgs []interface{}
	// 运行过程
	Course []string
	// 是否由HTTP模块调用引起的错误
	IsApi bool
	// 是否存在框架自身报错
	//  isApi		false
	//  isFrame		true
	//  ----		框架自身引起的报错
	//  isApi		true
	//  isFrame		true
	//  ----		由HTTP调用引起的报错,但同时调用了框架方法
	//  isApi		false
	//  isFrame		false
	//  ----		未知报错
	IsFrame bool
	// HTTP服务运行过程
	ApiCourse []string
	// 框架运行过程
	FrameCurse []string
}

type HttpF interface {
	Write(code int, content string, header ...map[string]string)
}
