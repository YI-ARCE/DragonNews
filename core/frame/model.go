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

type HttpRequest struct {
	Get  map[string]string
	Body string
}

type Error struct {
	Tag     string // 错误码
	Message string // 报错提示
	// 方法接收的内容,仅限于框架自身
	RequestArgs HttpRequest
	// 运行过程
	Course []string
	// HTTP服务运行过程
	ApiCourse []string
	// 框架运行过程
	FrameCourse []string
}

// Error 实现error接口
func (e *Error) Error() string {
	return e.Message
}

type HttpF interface {
	Body() []byte
	GetAll() map[string]string
}

// ErrorCode 错误码定义
const (
	// 系统错误
	SystemError = 500
	// 参数错误
	ParamError = 400
	// 权限错误
	AuthError = 403
	// 资源不存在
	NotFoundError = 404
	// 业务错误
	BusinessError = 1000
	// 数据库错误
	DatabaseError = 2000
	// 网络错误
	NetworkError = 3000
)
