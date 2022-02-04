package sqlmap

import (
	"NAS/src/config"
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	MyDb *gorm.DB
	db   *gorm.DB
	err  error
	mdb  config.Mysql
	rdb  config.Redis
)

func init() {
	mdb, rdb = config.GetConfig()
	log.Println("初始化数据库")
	MyDb = contGetDb()
	if MyDb == nil {
		log.Println("数据库初始化失败")
		os.Exit(-1)
	}
	err = MyDb.DB().Ping()
	if err != nil {
		log.Println("数据库初始化失败：", err.Error())
		os.Exit(-1)
		return
	}
}

func dbConn(User, Password, Host, Db string, Port int) *gorm.DB {
	connArgs := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", User, Password, Host, Port, Db)
	db, err = gorm.Open("mysql", connArgs)
	if err != nil {
		defer func(db *gorm.DB) {
			err = db.Close()
			if err != nil {
				return
			}
		}(db)
		return nil
	}
	db.SingularTable(true) //如果使用gorm来帮忙创建表时，这里填写false的话gorm会给表添加s后缀，填写true则不会
	db.LogMode(false)      //打印sql语句
	//开启连接池
	db.DB().SetMaxIdleConns(100)   //最大空闲连接
	db.DB().SetMaxOpenConns(10000) //最大连接数
	db.DB().SetConnMaxLifetime(30) //最大生存时间(s)

	return db
}

func contGetDb() (conn *gorm.DB) {
	for i := 0; i < 5; i++ {
		conn = dbConn(mdb.Username, mdb.Password, mdb.Host, mdb.Db, mdb.Port)
		if conn != nil {
			break
		}
		fmt.Println("本次未获取到mysql连接")
	}
	return conn
}
