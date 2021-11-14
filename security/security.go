package security

import (
	"errors"
	"fmt"
	"gin/pkg/settings"
	"strconv"
	"strings"
	"time"
)

var SessionAliveTime time.Duration

// rate limiting rules for APIs
var ApiLimitRules []LimitRule

// device register/login limit rule
var DeviceLimitRule LimitRule

// low risk => middle risk
var L2MRule LimitRule

// middle risk => high risk
var M2HRule LimitRule

func Setup() {

	// set session alive time
	sat := settings.SecuritySetting.SessionAliveTime
	if len(sat) < 2 {
		setupPanic()
	}
	size, err := strconv.Atoi(sat[:len(sat)-1])
	if err != nil {
		setupPanic()
	}
	unit := sat[len(sat)-1:]
	SessionAliveTime = time.Duration(size) * getDuration(unit)

	// resolve rate limiting rules
	ApiLimitRules = mapToRules(settings.SecuritySetting.ApiLimitRules)
	DeviceLimitRule = mapToRules(settings.SecuritySetting.DeviceLimitRule)[0]
	L2MRule = mapToRules(settings.SecuritySetting.L2MRule)[0]
	M2HRule = mapToRules(settings.SecuritySetting.M2HRule)[0]

}

// mapToRules "5/2m" => LimitRule{Count: 5, Window: 2 minutes}
func mapToRules(expression string) []LimitRule {
	many := strings.Split(expression, ",")
	var rules = make([]LimitRule, len(many))
	for i, m := range many {
		single := strings.Split(m, "/")
		if len(single) != 2 || len(single[1]) < 2 {
			setupPanic()
		}
		s0, s1 := single[0], single[1]
		count, err := strconv.Atoi(s0)
		if err != nil {
			setupPanic()
		}
		size, err := strconv.Atoi(s1[:len(s1)-1])
		if err != nil {
			setupPanic()
		}
		unit := s1[len(s1)-1:]
		window := time.Duration(size) * getDuration(unit)
		rules[i] = LimitRule{Count: count, Window: window}
	}
	return rules
}

func setupPanic() {
	panic(errors.New("security.Setup: Illegal expression"))
}

func getDuration(unit string) time.Duration {
	switch unit {
	case "s":
		return time.Second
	case "m":
		return time.Minute
	case "h":
		return time.Hour
	case "d":
		return time.Hour * 24
	default:
		panic(errors.New(fmt.Sprintf("security.Setup: Illegal time unit [%v]", unit)))
	}
}
