package main

//主方法

import (
	"fmt"
	"gin/models"
	"gin/pkg/gredis"
	"gin/pkg/settings"
	"gin/routers"
	"gin/security"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
	1.引入相应的库和包
	2.设置会用到的全局变量
	3.进行数据连接
	4.设置数据的结构体进行绑定操作
	5.解析并渲染网页模板操作
		a.定义模板
		b.解析模板
		c.渲染模板
	6.设计中间件
	7.撰写网页请求实现函数
		a.网页的GET请求
		b.网页的注册，登录的POST请求
        c.会话状态管理
    8.风控策略
        a.规则参数
        b.判断机制
        c.拦截层
	9.发布
*/

// init initialize server settings and connections
func init() {
	settings.Setup()
	models.Setup()
	gredis.Setup()
	security.Setup()
}

func main() {
	gin.SetMode(settings.ServerSetting.RunMode)

	routersInit := routers.InitRouter()
	readTimeout := settings.ServerSetting.ReadTimeout
	writeTimeout := settings.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", settings.ServerSetting.HttpPort)

	server := &http.Server{
		Addr:         endPoint,
		Handler:      routersInit,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	server.ListenAndServe()
}
