package settings

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

type app struct {
	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

var AppSetting = &app{}

type server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &server{}

type database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
	MaxIdleConn int
	MaxOpenConn int
}

var DatabaseSetting = &database{}

type redis struct {
	Addr        string
	Password    string
	DB          int
	IdleTimeout time.Duration
}

var RedisSetting = &redis{}

type security struct {
	ApplyCodeExpireTime int
	SessionAliveTime    string
	ApiLimitRules       string
	DeviceLimitRule     string
	L2MRule             string
	M2HRule             string
	TempBlockTime       int
}

var SecuritySetting = &security{}

var cfg *ini.File

// Setup initialize the configuration instance
func Setup() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("settings.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)
	mapTo("security", SecuritySetting)

	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
