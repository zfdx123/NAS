package main

import (
	"NAS/src/filesin"
	"NAS/src/middleware"
	"NAS/src/route"
	"log"
)

func init() {
	log.Println("初始化表")
	middleware.InitUser()
	filesin.InitFile()
}

func main() {
	var err error
	err = route.Route()
	if err != nil {
		log.Println("gin启动失败")
		return
	}
}
