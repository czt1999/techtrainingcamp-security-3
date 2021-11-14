package security

import (
	"fmt"
	"gin/pkg/gredis"
	"gin/pkg/settings"
	"log"
	"time"
)

// Normal：用户无异常行为
// LowRisk：进行滑块验证
// MidRisk：一段时间不能进行操作
// HighRisk：拒绝该用户后续所有请求
const (
	Normal   = 0
	LowRisk  = 1
	MidRisk  = 2
	HighRisk = 3

	BlockedCachePrefix = "block::"
)

type Env struct {
	IP       string `json:"ip"`
	DeviceID string `json:"device_id"`
}

// LimitRule 描述了在规定时间窗口内某一行为次数的上限
type LimitRule struct {
	Count  int
	Window time.Duration
}

//
// GetDecisionType 根据给定的env和userID，判断风险类型
//
// 项目根目录的文档指出，有两种情况会标记中等风险：
// 1）同一设备在规定时间内注册/登录用户数量达到规定值 CheckMidByDevice
// 2）同一IP/设备在规定时间内触发低等风险次数达到规定值 CheckMidByLow
//
// 从性能角度考虑，每个 Check 方法会先更新缓存中的窗口数据，再进行计数
// 所以，调用 CheckMidByLow 的前提是 CheckLow 为 true
// 调用 CheckHigh 的前提是 CheckMidByXX 其一为 true
//
func GetDecisionType(env Env, userID uint) int {
	isLow := CheckLow(env)
	isMid := CheckMidByDevice(env, userID)
	if isLow && CheckMidByLow(env) {
		isMid = true
	}
	if isMid {
		if CheckHigh(env) {
			return HighRisk
		}
		return MidRisk
	}
	if isLow {
		return LowRisk
	}
	return Normal
}

// CheckLow 检查是否达到低等风险
func CheckLow(env Env) bool {

	nowMs := time.Now().UnixNano() / 1000

	for i := 0; i < len(ApiLimitRules); i++ {
		// 第i条规则的缓存命名前缀为 limit[i]::
		windowSize := ApiLimitRules[i].Window.Milliseconds()
		cntIP, err := gredis.PutWindow(fmt.Sprintf("limit[%v]::%v", i, env.IP), nowMs, windowSize)
		if err != nil {
			log.Printf("PutWindow Error: %v\n", err)
		}
		cntDevice, err := gredis.PutWindow(fmt.Sprintf("limit[%v]::%v", i, env.DeviceID), nowMs, windowSize)
		if err != nil {
			log.Printf("PutWindow Error: %v\n", err)
		}
		if cntIP >= ApiLimitRules[i].Count || cntDevice >= ApiLimitRules[i].Count {
			log.Printf("CheckLow: %v is judged as LOW RISK\n", env)
			// clear window
			gredis.ClearWindow(fmt.Sprintf("limit[%v]::%v", i, env.IP), nowMs)
			gredis.ClearWindow(fmt.Sprintf("limit[%v]::%v", i, env.DeviceID), nowMs)
			return true
		}
	}

	return false
}

// CheckMidByLow 根据低等风险的触发次数，检查是否达到中等风险
// 该方法的调用必须在 CheckLow 判定为 true 之后
func CheckMidByLow(env Env) bool {

	nowMs := time.Now().UnixNano() / 1000
	windowSize := L2MRule.Window.Milliseconds()

	cntIP, err := gredis.PutWindow(fmt.Sprintf("l2m::%v", env.IP), nowMs, windowSize)
	if err != nil {
		log.Printf("PutWindow Error: %v\n", err)
	}
	cntDevice, err := gredis.PutWindow(fmt.Sprintf("l2m::%v", env.DeviceID), nowMs, windowSize)
	if err != nil {
		log.Printf("PutWindow Error: %v\n", err)
	}

	if cntIP >= L2MRule.Count || cntDevice >= L2MRule.Count {
		log.Printf("CheckMidByLow: %v is judged as MID RISK\n", env)
		// clear window
		gredis.ClearWindow(fmt.Sprintf("l2m::%v", env.IP), nowMs)
		gredis.ClearWindow(fmt.Sprintf("l2m::%v", env.DeviceID), nowMs)
		// set block information
		exp := time.Duration(settings.SecuritySetting.TempBlockTime) * time.Second
		_ = gredis.Set(BlockedCachePrefix+env.IP, "", exp)
		_ = gredis.Set(BlockedCachePrefix+env.DeviceID, "", exp)
		return true
	}
	return false
}

// CheckMidByDevice 根据设备注册/登录用户数量，检查是否达到中等风险
func CheckMidByDevice(env Env, userID uint) bool {
	if userID == 0 {
		return false
	}

	nowMs := time.Now().UnixNano() / 1000

	cnt, err := gredis.PutWindowWithValue(fmt.Sprintf("device-user::%v", env.DeviceID), userID, nowMs, DeviceLimitRule.Window.Milliseconds())
	if err != nil {
		log.Printf("PutWindow Error: %v\n", err)
	}

	if cnt >= DeviceLimitRule.Count {
		log.Printf("CheckMidByDevice: %v is judged as MID RISK\n", env)
		exp := time.Duration(settings.SecuritySetting.TempBlockTime) * time.Second
		_ = gredis.Set(BlockedCachePrefix+env.IP, "", exp)
		_ = gredis.Set(BlockedCachePrefix+env.DeviceID, "", exp)
		return true
	}
	return false
}

// CheckHigh 检查是否达到高等风险
// 该方法的调用必须在 CheckMidByLow 或者 CheckMidByDevice 判定为 true 之后
func CheckHigh(env Env) bool {

	nowMs := time.Now().UnixNano() / 1000
	windowSize := M2HRule.Window.Milliseconds()

	cntIP, err := gredis.PutWindow(fmt.Sprintf("m2h::%v", env.IP), nowMs, windowSize)
	if err != nil {
		log.Printf("PutWindow Error: %v\n", err)
	}
	cntDevice, err := gredis.PutWindow(fmt.Sprintf("m2h::%v", env.DeviceID), nowMs, windowSize)
	if err != nil {
		log.Printf("PutWindow Error: %v\n", err)
	}

	if cntIP >= M2HRule.Count || cntDevice >= M2HRule.Count {
		log.Printf("CheckHigh: %v is judged as HIGH RISK\n", env)
		// high level risk indicates both IP and deviceID will be blocked permanently
		_ = gredis.Set(BlockedCachePrefix+env.IP, "", 0)
		_ = gredis.Set(BlockedCachePrefix+env.DeviceID, "", 0)
		return true
	}
	return false
}
