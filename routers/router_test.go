package routers

import (
	"bytes"
	"encoding/json"
	"gin/models"
	"gin/pkg/app"
	"gin/pkg/gredis"
	"gin/pkg/settings"
	"gin/routers/api"
	"gin/security"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var router *gin.Engine

func init() {
	settings.Setup()
	models.Setup()
	gredis.Setup()
	security.Setup()
	router = InitRouter()
}

func TestCheckBlocked(t *testing.T) {
	req := api.RegisterRequest{}
	req.Username = "herry"
	req.PhoneNum = "13631199324"
	req.Password = "123456"
	req.Env = security.Env{
		IP:       "127.0.0.1",
		DeviceID: "xiaomi-v1",
	}
	s, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(s))
	req2.Header.Set("Content-Type", "application/json;charset=UTF-8")
	router.ServeHTTP(w, req2)
}

func TestRegister(t *testing.T) {
	req := api.RegisterRequest{}
	req.Username = "jerry"
	req.PhoneNum = "13631199324"
	req.Password = "123456"
	s, err := json.Marshal(req)

	//第一次注册，注册应该成功
	w := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(s))
	req2.Header.Set("Content-Type", "application/json;charset=UTF-8")
	router.ServeHTTP(w, req2)
	assert.Equal(t, http.StatusOK, w.Code)
	response := app.Response{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Equal(t, 0, response.Code)
	assert.Equal(t, "注册成功", response.Msg)

	//修改用户名，手机号不变
	req.Username = "tony"
	s, err = json.Marshal(req)
	w = httptest.NewRecorder()
	req2, _ = http.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(s))
	req2.Header.Set("Content-Type", "application/json;charset=UTF-8")
	router.ServeHTTP(w, req2)
	assert.Equal(t, http.StatusOK, w.Code)
	response = app.Response{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Equal(t, 1, response.Code)
	assert.Equal(t, "手机号已被注册", response.Msg)

	//修改手机号，用户名不变
	req.Username = "jerry"
	req.PhoneNum = "13631199325"
	s, err = json.Marshal(req)
	w = httptest.NewRecorder()
	req2, _ = http.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(s))
	req2.Header.Set("Content-Type", "application/json;charset=UTF-8")
	router.ServeHTTP(w, req2)
	assert.Equal(t, http.StatusOK, w.Code)
	response = app.Response{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Equal(t, 1, response.Code)
	assert.Equal(t, "用户名已被注册", response.Msg)

}

var SessionID string

func TestLoginByName(t *testing.T) {
	req := api.LoginByNameRequest{}
	req.Username = "jerry"
	req.Password = "123456"

	//密码用户名正确，登录成功
	s, err := json.Marshal(req)
	w := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodPost, "/api/login/name", bytes.NewReader(s))
	req2.Header.Set("Content-Type", "application/json;charset=UTF-8")
	router.ServeHTTP(w, req2)
	assert.Equal(t, http.StatusOK, w.Code)
	response := app.Response{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Equal(t, 0, response.Code)
	assert.Equal(t, "登录成功", response.Msg)

	m := response.Data.(map[string]interface{})
	for k, v := range m {
		switch value := v.(type) {
		case string:
			{
				if k == "session_id" {
					SessionID = value
				}
			}
		}
	}

	//用户名不变，密码不匹配
	req.Password = "1234"
	s, err = json.Marshal(req)
	w = httptest.NewRecorder()
	req2, _ = http.NewRequest(http.MethodPost, "/api/login/name", bytes.NewReader(s))
	req2.Header.Set("Content-Type", "application/json;charset=UTF-8")
	router.ServeHTTP(w, req2)
	assert.Equal(t, http.StatusOK, w.Code)
	response = app.Response{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Equal(t, 1, response.Code)
	assert.Equal(t, "用户名和密码不匹配", response.Msg)

	// 修改用户名，不修改密码
	req.Username = "tony"
	req.Password = "123456"
	s, err = json.Marshal(req)
	w = httptest.NewRecorder()
	req2, _ = http.NewRequest(http.MethodPost, "/api/login/name", bytes.NewReader(s))
	req2.Header.Set("Content-Type", "application/json;charset=UTF-8")
	router.ServeHTTP(w, req2)
	assert.Equal(t, http.StatusOK, w.Code)
	response = app.Response{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Equal(t, 1, response.Code)
	assert.Equal(t, "用户名和密码不匹配", response.Msg)

}

func TestLoginByPhone(t *testing.T) {
	req := api.LoginByPhoneRequest{}
	req.PhoneNum = "13631199324"

	// 先获取验证码
	reqcode := api.GetApplyCodeRequest{}
	reqcode.PhoneNum = "13631199324"
	s, err := json.Marshal(reqcode)
	w := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodPost, "/api/applycode", bytes.NewReader(s))
	req2.Header.Set("Content-Type", "application/json;charset=UTF-8")
	router.ServeHTTP(w, req2)
	assert.Equal(t, http.StatusOK, w.Code)
	response := app.Response{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	// 从data中解析验证码
	m := response.Data.(map[string]interface{})
	for k, v := range m {
		switch value := v.(type) {
		case string:
			{
				if k == "verify_code" {
					req.VerifyCode = value
				}
			}
		}
	}

	//手机号和验证码，登录成功
	s, err = json.Marshal(req)
	w = httptest.NewRecorder()
	req2, _ = http.NewRequest(http.MethodPost, "/api/login/phone", bytes.NewReader(s))
	req2.Header.Set("Content-Type", "application/json;charset=UTF-8")
	router.ServeHTTP(w, req2)
	assert.Equal(t, http.StatusOK, w.Code)
	response = app.Response{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Equal(t, 0, response.Code)
	assert.Equal(t, "登录成功", response.Msg)

	m = response.Data.(map[string]interface{})
	for k, v := range m {
		switch value := v.(type) {
		case string:
			{
				if k == "session_id" {
					SessionID = value
				}
			}
		}
	}
}

func TestLogout(t *testing.T) {

	req := api.LogoutRequest{}
	req.SessionID = SessionID

	// 登出账户

	req.ActionType = 1
	s, err := json.Marshal(req)
	w := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodPost, "/api/login/logout", bytes.NewReader(s))
	req2.Header.Set("Content-Type", "application/json;charset=UTF-8")
	router.ServeHTTP(w, req2)
	assert.Equal(t, http.StatusOK, w.Code)
	response := app.Response{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Equal(t, 0, response.Code)
	assert.Equal(t, "登出成功", response.Msg)

	// 重新登录才能注销，如果没有重新登录，session失效，无法注销
	req.ActionType = 2
	s, err = json.Marshal(req)
	w = httptest.NewRecorder()
	req2, _ = http.NewRequest(http.MethodPost, "/api/login/logout", bytes.NewReader(s))
	req2.Header.Set("Content-Type", "application/json;charset=UTF-8")
	router.ServeHTTP(w, req2)
	assert.Equal(t, http.StatusOK, w.Code)
	response = app.Response{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Equal(t, 1, response.Code)
	assert.Equal(t, "服务器内部错误", response.Msg)

	//重新登录
	reqLogin := api.LoginByNameRequest{}
	reqLogin.Username = "jerry"
	reqLogin.Password = "123456"

	//密码用户名正确，登录成功
	s, err = json.Marshal(reqLogin)
	w = httptest.NewRecorder()
	req2, _ = http.NewRequest(http.MethodPost, "/api/login/name", bytes.NewReader(s))
	req2.Header.Set("Content-Type", "application/json;charset=UTF-8")
	router.ServeHTTP(w, req2)
	assert.Equal(t, http.StatusOK, w.Code)
	response = app.Response{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Equal(t, 0, response.Code)
	assert.Equal(t, "登录成功", response.Msg)
	// 登出后sessionID失效，需要重新获得sessionID
	m := response.Data.(map[string]interface{})
	for k, v := range m {
		switch value := v.(type) {
		case string:
			{
				if k == "session_id" {
					SessionID = value
				}
			}
		}
	}

	// 注销账户
	req.SessionID = SessionID
	req.ActionType = 2
	s, err = json.Marshal(req)
	w = httptest.NewRecorder()
	req2, _ = http.NewRequest(http.MethodPost, "/api/login/logout", bytes.NewReader(s))
	req2.Header.Set("Content-Type", "application/json;charset=UTF-8")
	router.ServeHTTP(w, req2)
	assert.Equal(t, http.StatusOK, w.Code)
	response = app.Response{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Equal(t, 0, response.Code)
	assert.Equal(t, "注销成功", response.Msg)
}
