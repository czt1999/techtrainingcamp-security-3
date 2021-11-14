package api

import (
	"gin/models"
	"gin/pkg/app"
	"gin/pkg/gredis"
	"gin/pkg/util"
	"gin/security"
	"github.com/gin-gonic/gin/binding"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// @TODO 引入gin参数校验

type RegisterRequest struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	PhoneNum     string `json:"phone_number"`
	VerifyCode   string `json:"verify_code"`
	security.Env `json:"environment"`
}

const RegisterForbiddenCachePrefix = "registerx::"

// @Summary Add new user
// @Router /api/register [post]
func Register(c *gin.Context) {
	appG := app.Gin{C: c}

	req := RegisterRequest{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		fastFailIllegalArgs(&appG, req.Env)
		return
	}

	// verify phone number
	if VerifyMobileFormat(req.PhoneNum) == false {
		fastFailIllegalArgs(&appG, req.Env)
		return
	}

	// prevent frequent re-register
	if isExist, _ := gredis.Exist(RegisterForbiddenCachePrefix + req.PhoneNum); isExist {
		fastFail(&appG, req.Env, "请稍候再注册", nil)
		return
	}

	// verify apply code
	verifyCode, _ := gredis.Get(CodeCachePrefix + req.PhoneNum)
	if verifyCode != req.VerifyCode {
		fastFail(&appG, req.Env, "验证码错误", nil)
		return
	}

	// check duplication
	if isExist, err := models.ExistUserByName(req.Username); err != nil {
		fastFailInternalError(&appG, req.Env)
		return
	} else if isExist {
		fastFail(&appG, req.Env, "用户名已被注册", nil)
		return
	}

	if isExist, _, err := models.ExistUserByPhone(req.PhoneNum); err != nil {
		fastFailInternalError(&appG, req.Env)
		return
	} else if isExist {
		fastFail(&appG, req.Env, "手机号已被注册", nil)
		return
	}

	// encrypt password
	req.Password = util.EncodeMD5(req.Password)
	userID, err := models.AddUser(req.Username, req.Password, req.PhoneNum)
	if err != nil {
		fastFailInternalError(&appG, req.Env)
		return
	}

	// risk level decision
	decisionType := security.GetDecisionType(req.Env, userID)
	c.Set("decisionType", decisionType)

	if decisionType == security.Normal {

		// establish session
		sessionID, err := security.OpenSession(userID)
		if err != nil {
			appG.FailInternalError()
			return
		}

		appG.OK("注册成功", gin.H{
			"session_id":    sessionID,
			"expire_time":   "",
			"decision_type": decisionType,
		})

	} else {
		appG.Fail("注册失败", gin.H{
			"session_id":    "",
			"expire_time":   "",
			"decision_type": decisionType,
		})
	}
}

type LoginByNameRequest struct {
	Username     string `form:"username" json:"username"`
	Password     string `form:"password" json:"password"`
	security.Env `json:"environment"`
}

// @Summary User logins by name
// @Router /api/login/name [post]
func LoginByName(c *gin.Context) {
	appG := app.Gin{C: c}

	req := LoginByNameRequest{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		fastFailIllegalArgs(&appG, req.Env)
		return
	}
	log.Printf("Revcieve LoginByName request: %v", req)
	// encrypt password
	req.Password = util.EncodeMD5(req.Password)

	userID, err := models.CheckUser(req.Username, req.Password)
	if err != nil {
		fastFailInternalError(&appG, req.Env)
		return
	} else if userID == 0 {
		fastFail(&appG, req.Env, "用户名不存在或密码错误", nil)
		return
	}

	// risk level decision
	decisionType := security.GetDecisionType(req.Env, userID)
	c.Set("decisionType", decisionType)

	if decisionType == security.Normal {
		// establish session
		sessionID, err := security.OpenSession(userID)
		if err != nil {
			appG.FailInternalError()
			return
		}

		appG.OK("登录成功", gin.H{
			"session_id":    sessionID,
			"expire_time":   "",
			"decision_type": decisionType,
		})

	} else {
		appG.Fail("登录失败", gin.H{
			"session_id":    "",
			"expire_time":   "",
			"decision_type": decisionType,
		})
	}
}

type LoginByPhoneRequest struct {
	PhoneNum     string `json:"phone_number"`
	VerifyCode   string `json:"verify_code"`
	security.Env `json:"environment"`
}

// @Summary User logins by phone number
// @Router /api/login/phone [post]
func LoginByPhone(c *gin.Context) {
	appG := app.Gin{C: c}

	req := LoginByPhoneRequest{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		fastFailIllegalArgs(&appG, req.Env)
		return
	}

	// verify phone number
	if VerifyMobileFormat(req.PhoneNum) == false {
		fastFailIllegalArgs(&appG, req.Env)
		return
	}

	isExist, userID, err := models.ExistUserByPhone(req.PhoneNum)
	if err != nil {
		fastFailInternalError(&appG, req.Env)
		return
	} else if !isExist {
		fastFail(&appG, req.Env, "手机号未注册", nil)
		return
	}

	// get verify code from cache and do validation
	verifyCode, err := gredis.Get(CodeCachePrefix + req.PhoneNum)
	if err != nil {
		fastFailInternalError(&appG, req.Env)
		return
	}
	if verifyCode != req.VerifyCode {
		fastFail(&appG, req.Env, "验证码错误，请重新获取", nil)
		return
	}

	// risk level decision
	decisionType := security.GetDecisionType(req.Env, userID)
	c.Set("decisionType", decisionType)

	if decisionType == security.Normal {
		// establish session
		sessionID, err := security.OpenSession(userID)
		if err != nil {
			appG.FailInternalError()
			return
		}

		appG.OK("登录成功", gin.H{
			"session_id":    sessionID,
			"expire_time":   "",
			"decision_type": decisionType,
		})

	} else {
		appG.Fail("登录失败", gin.H{
			"session_id":    "",
			"expire_time":   "",
			"decision_type": decisionType,
		})
	}
}

type LogoutRequest struct {
	SessionID    string `json:"session_id"`
	ActionType   int    `json:"action_type"`
	security.Env `json:"environment"`
}

// @Summary User logouts
// @Router /api/register [post]
func Logout(c *gin.Context) {
	appG := app.Gin{C: c}

	req := LogoutRequest{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		fastFailIllegalArgs(&appG, req.Env)
		return
	}
	// get userID based on the session
	userID, err := security.GetUserID(req.SessionID)
	if userID == 0 {
		fastFail(&appG, req.Env, "用户未登录，操作非法", nil)
		return
	}

	switch req.ActionType {
	case 1:
		// close session
		if err = security.CloseSession(req.SessionID); err != nil {
			appG.FailInternalError()
			return
		}
		appG.OK("登出成功", nil)
	case 2:
		// close session
		if err = security.CloseSession(req.SessionID); err != nil {
			appG.FailInternalError()
			return
		}
		user, _ := models.GetUser(userID)
		// delete user
		if err = models.DeleteUser(userID); err != nil {
			appG.FailInternalError()
			return
		}
		// safety mechanism to prevent repeated register
		_ = gredis.Set(RegisterForbiddenCachePrefix+user.PhoneNum, "", 24*time.Hour)
		appG.OK("注销成功", nil)
	default:
		fastFail(&appG, req.Env, "操作不允许", nil)
	}
}

type GetUsernameRequest struct {
	SessionID    string `json:"session_id"`
	security.Env `json:"environment"`
}

func GetUsername(c *gin.Context) {
	appG := app.Gin{C: c}

	req := GetUsernameRequest{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		fastFailIllegalArgs(&appG, req.Env)
		return
	}
	log.Printf("Revcieve GetUsername request: %v", req)

	userID, err := security.GetUserID(req.SessionID)
	if err != nil {
		fastFailInternalError(&appG, req.Env)
		return
	} else if userID == 0 {
		fastFail(&appG, req.Env, "登录超时", nil)
		return
	}

	user, err := models.GetUser(userID)
	if err != nil {
		fastFailInternalError(&appG, req.Env)
		return
	} else if user.ID == 0 {
		fastFail(&appG, req.Env, "登录超时", nil)
		return
	}

	appG.OK("欢迎进入", gin.H{
		"username": user.Username,
	})

}

// fastFail, fastFailInternalError, fastFailIllegalArgs are designed for
// fast-failure against to bad request (i.e. with illegal args or unmatched
// SQL query result), which is commonly used in DDoS.
// After execute the fail action, it passes the env parameters to our security layer.
func fastFail(g *app.Gin, env security.Env, msg string, data interface{}) {
	g.Fail(msg, data)
	g.C.Set("decisionType", security.GetDecisionType(env, 0))
}

func fastFailInternalError(g *app.Gin, env security.Env) {
	g.FailInternalError()
	g.C.Set("decisionType", security.GetDecisionType(env, 0))
}

func fastFailIllegalArgs(g *app.Gin, env security.Env) {
	g.FailIllegalArgs()
	g.C.Set("decisionType", security.GetDecisionType(env, 0))
}
