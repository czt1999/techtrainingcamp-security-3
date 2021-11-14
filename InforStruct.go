package main



//这是一个数据结构体的类



import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

/*
	gorm.Model:gorm工具内数据模型
	Username:用户名称
	Phone：手机号
*/
type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(20);not null;"`
	Password string `gorm:"type:varchar(20);not null;"`
}

/*
	ID：数据库主键ID
	Username:用户名称
	Phone：手机号
	Ip：用户常用IP地址
*/
type UserInformation struct {
	Username string `gorm:"type:varchar(20);not null;"`
	Phone    string `gorm:"type:varchar(20);not null;"`
	Ip       string
}

/*
	页面表单提交数据
	Username:用户名称
	Password：密码
	Code:验证是否通过
*/
type UserLogin struct {
	Username string `gorm:"type:varchar(20);not null;"`
	Password string `gorm:"type:varchar(20);not null;"`
	Code     bool   ``
}

// Limiter 定义属性
type Limiter struct {
	// Redis client connection.
	rc *redis.Client
}

type VerifyCode struct {
	gorm.Model
	PhoneNumber string `gorm:"type:varchar(20);not null;"`
	VerifyCode  string `gorm:"type:varchar(20);not null;"`
	ExpireTime  int    `gorm:"type:int;not null;"`
}
