package database

import (
	"github.com/SSunSShine/QAsystem/conf"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var DB *gorm.DB

// 连接数据库
func init() {
	var err error

	DB, err = gorm.Open(conf.Config().DB.Driver, conf.Config().DB.Addr)

	if err != nil {
		log.Println(err)
		panic("failed to connect database !")
	}

	// 设置连接池连接数
	DB.DB().SetMaxOpenConns(150)
	DB.DB().SetMaxIdleConns(100)

	DB.SingularTable(true)
}