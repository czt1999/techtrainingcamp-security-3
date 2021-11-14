package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	OK   = 0
	FAIL = 1
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response wrap gin.JSON
func (g *Gin) Response(code int, message string, data interface{}) {
	g.C.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  message,
		Data: data,
	})
	return
}

func (g *Gin) OK(message string, data interface{}) {
	g.Response(OK, message, data)
}

func (g *Gin) Fail(message string, data interface{}) {
	g.Response(FAIL, message, data)
}

func (g *Gin) FailIllegalArgs() {
	g.Fail("请求参数错误", nil)
}

func (g *Gin) FailInternalError() {
	g.Fail("服务器内部错误", nil)
}
