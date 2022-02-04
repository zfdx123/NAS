package config

import (
	"github.com/Unknwon/goconfig"
	"log"
	"os"
)

type (
	Mysql struct {
		Username string
		Password string
		Host     string
		Port     int
		Db       string
	}
	Redis struct {
		Password string
		Addr     string
		Db1      int
		Db2      int
	}
	Web struct {
		TLS    bool
		Addr   string
		TLSCrt string
		TLSKey string
	}
)

func GetConfig() (Mysql, Redis) {
	var mysql Mysql
	var redis Redis
	cfg, err := goconfig.LoadConfigFile("./file/conf.ini")
	if err != nil {
		log.Println("配置文件处理失败", err.Error())
		os.Exit(-1)
	}
	mysql = Mysql{
		Username: cfg.MustValue("Mysql", "Username", "root"),
		Password: cfg.MustValue("Mysql", "Password", "root"),
		Host:     cfg.MustValue("Mysql", "Host", "127.0.0.1"),
		Port:     cfg.MustInt("Mysql", "Port", 3306),
		Db:       cfg.MustValue("Mysql", "DB"),
	}
	redis = Redis{
		Password: cfg.MustValue("Redis", "Password", ""),
		Addr:     cfg.MustValue("Redis", "Addr", "127.0.0.1:6379"),
		Db1:      cfg.MustInt("Redis", "DB1", 0),
		Db2:      cfg.MustInt("Redis", "DB2", 1),
	}
	return mysql, redis
}

func GetWebConfig() Web {
	var webC Web
	cfg, err := goconfig.LoadConfigFile("./file/conf.ini")
	if err != nil {
		log.Println("配置文件处理失败", err.Error())
		os.Exit(-1)
	}
	webC = Web{
		TLS:    cfg.MustBool("Web", "TLS", false),
		Addr:   cfg.MustValue("Web", "Addr", ":8080"),
		TLSCrt: cfg.MustValue("Web", "TLSCrt", "./file/web.pem"),
		TLSKey: cfg.MustValue("Web", "TLSKey", "./file/web.crt"),
	}
	return webC
}
