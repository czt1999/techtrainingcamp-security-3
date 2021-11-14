package api

import (
	"fmt"
	"gin/pkg/app"
	"gin/pkg/gredis"
	"gin/pkg/settings"
	"gin/security"
	"github.com/gin-gonic/gin/binding"
	"log"
	"math/rand"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

type GetApplyCodeRequest struct {
	PhoneNum     string `json:"phone_number"`
	security.Env `json:"environment"`
}

const CodeCachePrefix = "applycode::"

// @Summary Get apply code based on phone number
// @Router /api/applycode [post]
func GetApplyCode(c *gin.Context) {
	appG := app.Gin{C: c}

	req := GetApplyCodeRequest{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		fastFailIllegalArgs(&appG, req.Env)
		return
	}

	// verify phone number
	if VerifyMobileFormat(req.PhoneNum) == false {
		fastFailIllegalArgs(&appG, req.Env)
		return
	}

	// get dicision type
	decisionType := security.GetDecisionType(req.Env, 0)

	if decisionType == security.Normal {
		// 生成六位随机数
		verifyCode := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
		expireTime := settings.SecuritySetting.ApplyCodeExpireTime

		// cache the verifyCode with expireTime
		if err := gredis.Set(CodeCachePrefix+req.PhoneNum, verifyCode, time.Duration(expireTime)*time.Second); err != nil {
			appG.FailInternalError()
			return
		}

		appG.OK("获取验证码成功", gin.H{
			"verify_code":   verifyCode,
			"expire_time":   expireTime,
			"decision_type": decisionType,
		})
		log.Printf("Send verify code [%v] to [%v]\n", verifyCode, req.PhoneNum)
	} else {

		appG.Fail("获取验证码失败", gin.H{
			"verify_code":   "",
			"expire_time":   -1,
			"decision_type": decisionType,
		})
	}
}

const verifyMobileRegular = "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"

// VerifyMobileFormat mobile verify
// 正则表达式检验手机号是否合法
func VerifyMobileFormat(mobileNum string) bool {
	reg := regexp.MustCompile(verifyMobileRegular)
	return reg.MatchString(mobileNum)
}
