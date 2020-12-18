#bingo-mall
Simple Rest API using Gin web framework, Gorm as ORM 、Logrus as log middleware ,also using MVC pattern 
基于 Gin和Gorm 搭建的一个Web后台管理系统的框架。以gin-first进行了封装开发的一套在线商城系统，包括了后台，还有小程序前端，vue的后端，现在还是比较粗糙。部分功能没有实现，后续会继续投入开发。
特点：
- 内置很多web常用的正则校验规则，以及正则校验函数(helper/regex.go)
- 基于JWT的API接口token认证
- 基于Logrus实现日志分类输出
- 基于gin-swagger实现，api接口文档输出
- 实现跨域Cors
- 配置好了完整的dockerfile，可以直接运行将项目编译成docker镜像
在项目根目录下执行： `swag init`

Swagger效果：http://127.0.0.1:8088/swagger/index.html (http://timesweb.com/swagger/index.html)
bwh: /var/data/gopath/bin/swag init

![image](https://github.com/YinYongTao/bingomall/blob/master/view/API%E9%A2%84%E8%A7%88%E6%95%88%E6%9E%9C.png)

model中可以这样子用，用来存多个url
//ImageUrl pq.StringArray `gorm:"type:varchar(5000)[];" form:"imageUrl[]" binding:"required" json:"imageUrl"`
ImageUrl string `gorm:"type:varchar(5000);" form:"imageUrl" binding:"required" json:"imageUrl"`

编辑器：http://idea.lanyus.com/

###如何启动
1. 把constant/wechat.go.example改成constant/wechat.go，把里面的相关参数改成自己的
2. 把conf/datasource.yml里面的参数改成自己的参数
3. 然后在bingomall的根目录下运行go run main.go就可以跑起来了
4. 相关api文档见swag部分
