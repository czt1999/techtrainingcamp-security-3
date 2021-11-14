package routers

import (
	"gin/routers/api"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())

	// API
	apiG := r.Group("/api")
	apiG.POST("/register", api.Register)
	apiG.POST("/login/name", api.LoginByName)
	apiG.POST("/login/phone", api.LoginByPhone)
	apiG.POST("/logout", api.Logout)
	apiG.POST("/applycode", api.GetApplyCode)
	apiG.POST("/user/name", api.GetUsername)

	// interceptor
	apiG.Use(AddBlocked)
	apiG.Use(CheckBlocked)

	// static resource (css, js)
	r.Static("pages/css", "./template/pages/css")
	r.Static("pages/js", "./template/pages/js")

	// resolve HTML templates
	r.LoadHTMLGlob("template/*.html")

	r.GET("/", func(c *gin.Context) {
		log.Printf("GET request from: %v\n", GetIp(c))
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.GET("/index", func(c *gin.Context) {
		log.Printf("GET request from: %v\n", GetIp(c))
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.GET("/main", func(c *gin.Context) {
		log.Printf("GET request from: %v\n", GetIp(c))
		c.HTML(http.StatusOK, "main.html", gin.H{})
	})

	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "error.html", nil)
	})
	return r
}

//获取登录页面的访问IP地址
func GetIp(c *gin.Context) string {
	reqIP := c.ClientIP()
	if reqIP == "::1" {
		reqIP = "127.0.0.1"
	}
	return reqIP
}
