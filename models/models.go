package models

import (
	"fmt"
	"gin/pkg/settings"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

// Setup initializes the database instance
func Setup() {

	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		settings.DatabaseSetting.User,
		settings.DatabaseSetting.Password,
		settings.DatabaseSetting.Host,
		settings.DatabaseSetting.Name)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// if necessary, set your own config
	// &gorm.Config{
	//	NamingStrategy: schema.NamingStrategy{
	//		TablePrefix:   settings.DatabaseSetting.TablePrefix, // 表名前缀，`User`表为`{prefix}_users`
	//		SingularTable: true,                                 // 使用单数表名，启用该选项后，`User` 表将是`user`
	//	},
	// }

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	sqlDB.SetMaxIdleConns(settings.DatabaseSetting.MaxIdleConn)
	sqlDB.SetMaxOpenConns(settings.DatabaseSetting.MaxOpenConn)

	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}
}
