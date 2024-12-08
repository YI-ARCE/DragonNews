更新记录
====================

* 1.0.0
    * 减少了繁杂配置项
    * 基础路由复刻
    * 调用简单
    * 适配Db,使用方法基本相同
    * 抛弃view视图支持

* 1.0.1
    * 修复出现并发后请求数据被覆盖的问题(更新参考[README.md](https://github.com/YI-ARCE/DragonNews/blob/main/README.md))

* 1.0.2
    * 增加cache缓存模块(目前仅支持redis)(更新参考[README.md](https://github.com/YI-ARCE/DragonNews/blob/main/README.md))

* 1.0.3
    * 增加log日志输出,支持大部分数据以json或文本型输出,输出目录为根目录->当前年月->当前日.txt
    * 优化orm并增加事物处理
    * 增加微信及微信支付相关SDK
    * 集成支付宝及阿里云SDK
    * 增加定时器模块
    * 增加unipushSDK,目前测试可正常使用
    * 增加配置获取,仅支持yaml,并将框架配置文件移至根目录
    * 增加curl模块

* 2.0
  * 沉淀之后的完全更新,抛弃纯TP思路,在基础上增加go的特性思想
  * readme更新更新参考[README.md](https://github.com/YI-ARCE/DragonNews/blob/main/README.md)
  
