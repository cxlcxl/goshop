package connect_mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func ConnectMysql(dsn string) (db *gorm.DB, err error) {
	db, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Fatal("数据库连接失败", err)
	}
	_db, err := db.DB()
	if err != nil {
		log.Fatal("数据库连接失败 1", err)
	}
	_db.SetMaxIdleConns(10)
	_db.SetMaxOpenConns(100)
	return
}
