package routers

import (
	"gin/pkg/gredis"
	"gin/pkg/settings"
	"gin/security"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type BaseRequest struct {
	security.Env `json:"environment"`
}

const BlockedCachePrefix = "block::"

// CheckBlocked check risk level of the request and determine whether it should be blocked
func CheckBlocked(c *gin.Context) {
	req := BaseRequest{}
	if err := c.BindJSON(&req); err != nil {
		// request without Env prameters will be intercepted directly
		c.Abort()
		return
	}

	if req.IP == "" {
		req.IP = GetIp(c)
	}

	// get block information from cache
	e, err := gredis.Exist(BlockedCachePrefix + req.IP)
	if err != nil || e {
		c.Abort()
		return
	}
	e, err = gredis.Exist(BlockedCachePrefix + req.DeviceID)
	if err != nil || e {
		c.Abort()
		return
	}

	// pass
	c.Next()
}

func AddBlocked(c *gin.Context) {

	c.Next()

	req := BaseRequest{}
	_ = c.BindJSON(&req)

	// after handling the request, we check decisionType and generate block information
	decisionType := c.GetInt("decisionType")
	switch decisionType {
	case security.MidRisk:
		exp := time.Duration(settings.SecuritySetting.TempBlockTime) * time.Second
		if err := gredis.Set(BlockedCachePrefix+req.IP, "", exp); err != nil {
			log.Fatalln(err)
			return
		}
		if err := gredis.Set(BlockedCachePrefix+req.DeviceID, "", exp); err != nil {
			log.Fatalln(err)
			return
		}
	case security.HighRisk:
		// high level risk indicates both IP and deviceID will be blocked permanently
		if err := gredis.Set(BlockedCachePrefix+req.IP, "", 0); err != nil {
			log.Fatalln(err)
			return
		}
		if err := gredis.Set(BlockedCachePrefix+req.DeviceID, "", 0); err != nil {
			log.Fatalln(err)
			return
		}
	default:
	}

}
