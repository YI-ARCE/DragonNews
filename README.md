DragonNews
===============


基于ThinkPHP5.0复刻而来

 + 减少了繁杂配置项
 + 基础路由复刻
 + 调用简单
 + 适配Db,使用方法基本相同
 + 抛弃view视图支持

> 运行环境最好在Go1.15往上。

## 目录结构

初始的工程目录结构如下：

~~~
├─application           应用目录
│  ├─packages           包目录
│  │  ├─index           模块目录
│  │    └─index.go      模块文件
│  │  └─ ...            更多模块目录
│  │
│  ├─config.yaml        整体配置文件
│  ├─go.mod             Go依赖文件
│  ├─go.sum             Go依赖文件
│  └─Route.go           路由注册文件
│
├─dragonnews            框架系统目录
│  ├─base               
│  ├─config             配置加载模块
│  ├─http               http启动模块
│  ├─orm                orm模块
│  ├─packages           
│  ├─pkg                pkg包
│  ├─reply              请求解析模块
│  ├─route              路由模块
│  ├─go.mod             Go依赖文件
│  ├─go.sum             Go依赖文件
│  └─Start.go           启动文件
│
├─go.mod                Go依赖文件
├─go.sum                Go依赖文件
├─main.go               调试启动文件
~~~

config.yaml配置参考
---------------
~~~
#orm配置
sql :
  #数据库类型,目前仅支持mysql类型
  type : 'mysql'
  #连接地址
  host : '127.0.0.1'
  #连接端口
  port : '3306'
  #连接库名
  database : 'dn_database'
  #用户名
  username : 'root'
  #密码
  password : 'root'

#框架配置
dragonnews :
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

    浏览器访问：127.0.0.1:8080/index
    出现"Hello World"即为成功运行

    //编译
    go build main.go
~~~

路由加载
---------------
* application/Route.go
~~~
    //get路由
    route.Get("/index",Index.Index)

    //post路由
    route.Post("/index",Index.Index)

    //无区分路由,当该路由存在于get或post时,
    //会根据当前方式加载对应方法,需要注意重名
    route.rule("/index",Index.Index)
~~~

> 以下例子目录: packages/index/index.go

数据库类使用
---------------
* 查询一条记录
~~~
    func Index(reply reply.Reply) {
    	result,_ := Db.Table("dn_table").Where("dn_id = 4").Find()
    	fmt.Print(result)
    }

    //结果返回一个map,为空时返回空map
    map[dn_id:3 dn_name:龙讯框架]
~~~

* 查询多条记录
~~~
    func Index(reply reply.Reply) {
    	result,_ := Db.Table("dn_table").Where("dn_id != 4").Select()
    	fmt.Print(result)
    }

    //结果返回查出的多个map,为空时返回空map
    map[0:map[dn_id:1 dn_name:龙讯框架] 1:map[dn_id:2 dn_name:龙讯科技,掌握核爆之力]]
~~~

* 更新
~~~
    func Index(reply reply.Reply) {
        up := map[string]string{
            "dn_name":"龙讯时代",
        }
    	result,_ := Db.Table("dn_table").Where("dn_id != 4").Update(up)
    	fmt.Print(result)
    }

    //结果返回更新的条数,为零时只有意外情况未执行或执行了没有更新的状态,根据实际情况判断
    2
~~~

* 插入
~~~
    func Index(reply reply.Reply) {
        ins := map[string]string{
            "dn_name":"龙讯力量",
        }
    	result,_ := Db.Table("dn_table").Insert(up)
    	fmt.Print(result)
        //返回本次插入的主键ID
        fmt.Print(Db.GetLastId())
    }

    //结果返回真,返回假则插入失败
    true
    //为0时执行失败
    5
~~~

* 删除
~~~
    func Index(reply reply.Reply) {
    	result,_ := Db.Table("dn_table").Where("dn_id = 5").Delete()
    	fmt.Print(result)
    }

    //结果返回删除的条数,判断条件与更新相同
    1
~~~

* 生成SQL语句只需要调用时跟上Fet()即可,注意获取结果时,第一个参数为结果,当生成时获取第二个结果
~~~
    func Index(reply reply.Reply) {
    	_,str := Db.Table("dn_table").Where("dn_id = 1").Fet().Find()
    	fmt.Println(str)
    }

    //结果返回生成好的SQL字符串
    SELECT * FROM `dn_table` WHERE (dn_id = 1) LIMIT 1
~~~

请求数据获取
----------------
* 获取GET数据, 请求uri: /index?dn_name=龙讯传说&db_id=1
~~~
    func Index(reply reply.Reply) {
        result := reply.Request.Get
        fmt.Print(result)
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
    func Index(reply reply.Reply) {
        result := reply.Request.Body
        fmt.Print(result)
    }
    
    //结果返回map[string]interface{}类型,使用需要断言
    map[dn_name:龙讯科技 dn_url:search.yiarce.cn]
~~~
~~~
    //如果想自己接收数据
    func Index(reply reply.Reply) {
    	data := make(map[string]string)
    	err := json.Unmarshal(reply.Request.BodyByte,&data)
    	if err != nil {
            fmt.Println(err)
        }
    	fmt.Print(data)
    }
    
    //结果,使用不需要断言
    map[dn_name:龙讯科技 dn_url:search.yiarce.cn]
~~~

* 获取POST表单,兼容form-data和x-www-form-urlencoded,都使用该方式获取即可
~~~
    func Index(reply reply.Reply) {
    	data := reply.Request.Post
        fmt.Println(data)
    }

    //结果返回map
    map[dn_name:龙讯科技 dn_url:search.yiarce.cn]
~~~

* 获取文件,上传键名为dn_file
~~~
    func Index(reply reply.Reply) {
    	data := reply.Request.File["dn_file"]
        fmt.Println(data)
    }

    //结果返回文件二进制数据 []byte类型
    [137 80 78 71 13 10 26 10 0 0 0 13 73 72 68 82 ...]
~~~

返回数据
--------------
* 设置header
~~~
    func Index(reply reply.Reply) {
    	reply.SetHeader("Content-Length","1000")
    }
~~~

* 结束请求并返回数据
~~~
    func Index(reply reply.Reply) {
    	data := map[string]string{
            "dn_name":"龙讯科技",
            "dn_url":"search.yiarce.cn/",
        }
    	reply.Return(200,data)
    }

    //Postman获取结果为json
    {
        "dn_name": "龙讯科技",
        "dn_url": "search.yiarce.cn/"
    }
~~~

~~~
    调用reply.Return时,参数一为浏览器状态码,参数二为返回的数据,参数三为返回的格式,可以不填
    会根据数据类型自动转换,未识别时会自动转换为json格式
    目前可识别的为string,byte两个类型,若为struct或map时会自动调用json.Marshal()返回
    //不填参数三时
    reply.Return(200,data)
    //填参数三
    reply.Return(200,data,reply.Ct.Json)
    如果想自己处理可以使用以下方法


    reply.W.Header("Content-Type","application/json")
    //如果想要返回别的状态码时
    reply.W.WriteHeader(500)


    _, _ = reply.W.Write([]byte(`{"dn_name":"龙讯科技","dn_url":"search.yiarce.cn"}`))
    
    //返回结果
    {
        "dn_name": "龙讯科技",
        "dn_url": "search.yiarce.cn/"
    }
~~~

Cache
--------------
* 配置
~~~

~~~
> [更新日志](https://github.com/YI-ARCE/DragonNews/blob/main/UPDATE.md)


## 版权信息

DragonNews遵循MIT开源协议发布，并提供免费使用。
