DragonNews
===============
> 运行环境在go 1.20 以上最佳

## 目录结构

初始的工程目录结构如下：

~~~
├─app                   应用目录
│  └─{controller}       路由分级目录
│     ├─api.go          定义请求接口文件
│     ├─model.go        定义接口模型
│     └─{controller}    子集路由,结构相同
│
├─config                配置目录
│  ├─frame.yaml         框架配置        
│  ├─database.yaml      数据库配置
│  └─log                日志配置
├─core                  框架提供的组件
│  ├─cache              缓存
│  ├─curl               请求
│  ├─date               日期
│  ├─dhttp              请求及路由处理
│  ├─encrypt            加密
│  ├─file               文件读写
│  ├─frame              运行时的组件
│  ├─log                日志
│  ├─timing             定时器
│  └─yorm               框架提供的sql-orm
│
├─go.mod                Go依赖文件
├─go.work               框架注册定义
├─main.go               启动入口
~~~

config.yaml配置参考
---------------
~~~
#数据库配置
database :
  # 数据库类型,目前仅支持mysql或兼容mysql语法的数据库
  type : 'mysql'
  # 连接地址
  host : '127.0.0.1'
  # 连接端口
  port : '3306'
  # 连接库名
  database : 'dn_database'
  # 用户名
  username : 'root'
  # 密码
  password : 'root'

#框架配置
frame :
  #请求解析(目前未实装)
  request :
    #是否开启加密
    enable : false
    #开启后的解密方式,支持AES和RSA
    secretMode : 'AES'
    #公钥,AES只需填写16位公钥即可
    public : '1234567890123456'
    #私钥
    private : ''

  #响应头
  response:
    #错误通知,开启后会可替换指定位置的文本内容为错误信息
    errorNotice: true
    #默认响应类型
    returnType: 'json'
    #当路由解析不到地址或出现错误时的浏览器状态码
    statusCode: 500
    #返回信息,可以使用替换标识[errorMsg]对文本替换报错内容
    errorData: '{"code":"50000","message":"[errorMsg]","data":{}}'

  #服务启动
  server :
    #监听地址
    host : '0.0.0.0'
    #监听端口
    port : '8080'
~~~



运行项目
---------------
* 运行前如要修改先配置config.yaml,最好配置一下go module
~~~
    //如果出现依赖项未加载未解析等,在项目根中执行以下命令
    go mod tidy

    //运行项目在项目根中执行以下命令
    go run main.go

    浏览器访问：127.0.0.1:8080/index/hello
    出现"Hello World"即为成功运行

    //编译
    go build main.go
~~~

路由加载
---------------
* router/register.go
~~~
    //get路由
    dhttp.Get("/index",Index.Index)

    //post路由
    dhttp.Post("/index",Index.Index)

    //无区分路由,当该路由存在于get或post时,默认get
    //会根据当前方式加载对应方法,需要注意重名
    dhttp.rule("/index",Index.Index)
~~~

> 以下例子目录: app/login、app/user

数据库类使用
---------------
* 查询一条记录
~~~
func Phone(d *dhttp.Dn) {
    result := d.Table(`xxx`).Where(`a = 2`).Find()
    if result.Err() != nil {
        d.Json(result.Err().Error(), 500)
    }
    m := sakuraPost.Structure{}
    result.Format(&m)
    d.Json(m)
}

    //结果返回一个map[string]string,为空时返回空map
    map[dn_id:3 dn_name:龙讯框架]
~~~

* 生成SQL语句只需要调用时跟上Fet()即可,注意获取结果时,第一个参数为结果,当生成时获取第二个结果
~~~
    func Index(d *dhttp.Dn) {
    	result := d.Table("dn_table").Where("dn_id = 1").Fet(true).Find()
    	frame.Println(result.Sql())
    }

    //结果返回生成好的SQL字符串
    SELECT * FROM `dn_table` WHERE (dn_id = 1) LIMIT 1
~~~

请求数据获取
----------------
* 获取GET数据, 请求uri: /index?dn_name=龙讯传说&db_id=1
~~~
    func GetKey(d *dhttp.Dn) {
        d.Json(d.GetAll())
    }
    
    //结果返回map
    map[db_id:1 dn_name:龙讯传说]
~~~

* 获取POST数据
~~~
    //请求json
    {
        "dn_name":"龙讯科技",
        "dn_url":"search.yiarce.cn"
    }
~~~
~~~
func GetJson(d *dhttp.Dn) {
    // 如果提交数据复杂,请使用Format
    //data := map[string]interface{}{}
    //d.BodyFormat(data, `服务器错误`)
    // ------------------------
    // 返回map[string]string类型
    frame.Println(d.Data())
    d.Json(d.Data())
}
    
    //结果返回map[string]string类型
    map[dn_name:龙讯科技 dn_url:dragon-news.yiarce.cn]
~~~
~~~
    //如果想自己接收数据
    func Index(d *dhttp.Dn) {
    	data := make(map[string]string)
    	err := json.Unmarshal(d.Body(),&data)
    	if err != nil {
            frame.Println(err)
        }
    	frame.Println(data)
    }
    
    //结果,使用不需要断言
    map[dn_name:龙讯科技 dn_url:dragon-news.yiarce.cn]
~~~

* 获取POST表单,兼容form-data和x-www-form-urlencoded,都使用该方式获取即可
~~~
    func Index(d *dhttp.Dn) {
    	name := d.Post(`dn_name`)
    	// 如果表单中存在文件流,文件流会在后续更改结构体,方便开发者处理
        image := d.File(`dn_image`)
        frame.Println(name)
        frame.Println(image)
    }

    //结果
    name->  龙讯框架
    image-> [137 80 78 71 13 10 26 10 0 0 0 13 73 72 68 82 ...]
~~~

返回数据
--------------
* 设置header
~~~
    func Index(d *dhttp.Dn) {
    	d.SetHeader("Content-Length","1000")
    }
~~~


自动路由生成
=======

* 以下是一个标准接口方法定义
~~~
// Phone 手机号登录
//
//  -param PhoneModel
//  -method post
func Phone(d *dhttp.Dn) {
    result := d.Table(`xxx`).Where(`c = 1`).Find()
    if result.Err() != nil {
        d.Json(result.Err().Error(), 500)
    }
    d.Json(result.Result(), 500)
}

// 结构体内容如下
// PhoneModel 手机号登录
type PhoneModel struct {
    Type     int    `json:"type"`               // 登录类型,1->密码,2->验证码
    Phone    string `json:"phone"`              // 手机号
    Code     string `json:"code,omitempty"`     // 验证码
    Password string `json:"password,omitempty"` // 密码
}
~~~

* param定义一个请求结构体,在当前目录中model中定义该结构体
* param允许为空或不声明,为空时自动寻找model中方法名+Model的结构体
* 如果未寻到符合的结构体则默认该请求未定义请求参数,不参与生成接口文档
* method定义请求方式,为空或不声明默认get

数据库表结构生成
=======

* 准备生成时请配置好database连接参数,以下是一个标准的数据库表
* 对数据库字段命名注释规则较高,也是为了全员适用一种注释规范
~~~

~~~

* param定义一个请求结构体,在当前目录中model中定义该结构体


## 版权信息

DragonNews遵循MIT开源协议发布，并提供免费使用。
