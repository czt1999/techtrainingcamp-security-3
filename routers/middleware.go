package routers

import (
	"gin/pkg/gredis"
	"gin/security"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseRequest struct {
	security.Env `json:"environment"`
}

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
	e, err := gredis.Exist(security.BlockedCachePrefix + req.IP)
	if err != nil || e {
		c.JSON(http.StatusTooManyRequests, "你已被禁止访问")
		c.Abort()
		return
	}
	e, err = gredis.Exist(security.BlockedCachePrefix + req.DeviceID)
	if err != nil || e {
		c.JSON(http.StatusTooManyRequests, "你已被禁止访问")
		c.Abort()
		return
	}

	// pass
	c.Next()
}