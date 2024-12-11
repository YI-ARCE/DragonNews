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

config 配置参考
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

- 构建表结构时对注释有一定的要求,为了风格统一防止查看注释各有个的写法之外也是配合框架做辅助工作<br>
- `[dn:time]` 该标签注释用于字段是时间戳的值
- `[dn:encrypt]` 该标签注释用于字段是加密的值
  标签用于框架根目录下的`sync_table.exe`程序使用,若未正确填写可能导致识别异常<br>
  该程序将会同步当前配置的数据库表转化为go结构体使用<br>
  当数据库表有新改动时请运行一次更新框架的表结构并提交git<br>
~~~
// 表构建不要求所有字段都需要有备注,但表名备注必须要有
// 对于字段存在多个意义的情况下,应该用value->remark,value->remark的格式注释
// 请严格使用英文逗号做分隔符
// 个别字段可能存在多类型分支的状态表达
// 如: 
// A,1->状态1,2->状态2
// B,1->状态1,2->状态2
// 此时可以通过竖线符 | 进行拼接
// A,1->状态1,2->状态2|B,1->状态1,2->状态2

// 例如create_time,update_time这些通用化字段可以不需要中文注释
// 但要求必须加上[dn:time]标签注释

// mobile这里被要求使用加密方式解析数据

create table user_info
(
    xx_id       int auto_increment  primary key,
    status      tinyint      default 0  null comment '状态,0->初始状态,1->开启,4->删除',
    mobile      varchar(128)      default 0  null comment '手机号码,[dn:encrypt]',
    create_time int          default 0  not null '创建时间,[dn:time]',
    update_time int          default 0  not null '更新时间,[dn:time]'
) comment 'XX表';
~~~


- 日志存储
  >   在使用`success`,`error`时会自动调用日志存储组件将本次日志写入库中<br>
  如果被异常拦截无法触发时,frame会代替写入日志并存储本次的异常信息
  ----
  | ID | 用户ID | 记载类型 | 日志ID    | 请求数据      | 响应数据                 | 请求时间       | 响应时间       |
        |:---|:-----|:-----|:--------|:----------|:---------------------|:-----------|:-----------|
  | 1  | 271  | 2    | mas312y | {"a":"b"} | {"msg":"success"...} | 1709602016 | 1709602017 |
  ````
    func GetKey(d *dhttp.Dn) {
        d.Json(`{}`,200)
    }
  ````



## 版权信息

DragonNews遵循MIT开源协议发布，并提供免费使用。
