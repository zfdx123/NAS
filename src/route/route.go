package route

import (
	"NAS/src/config"
	"NAS/src/filesin"
	md "NAS/src/middleware"
	"NAS/src/sysproc"
	"NAS/src/ws"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func Route() error {
	webConfig := config.GetWebConfig()
	//r := gin.Default()
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), gzip.Gzip(gzip.DefaultCompression))
	err := r.SetTrustedProxies(nil)
	if err != nil {
		return err
	}

	v1 := r.Group("/T")
	//v1.Use(md.JWTAuth())
	{
		v1.POST("/login", md.Login)
		v1.POST("/new", md.New)
	}

	user := r.Group("/User")
	{
		user.POST("/rest", md.RestPass)
		user.GET("/Logout", md.Logout)
	}

	sys := r.Group("/System")
	sys.Use(md.JWTAuth(), SetUp())
	{

		sys.GET("/sys-info", sysproc.GetSysInfo)
		sys.POST("/proc", sysproc.Shutdown)
		sys.GET("/ws", ws.Ws)
	}
	file := r.Group("/File")
	file.Use(md.JWTAuth(), SetUp()) // md.JWTUpd(),
	{
		file.POST("/Uploads", filesin.UpdateLoad)                //上传文件确认接口
		file.POST("/Upload", filesin.Upload)                     //文件上传接口
		file.GET("/Download/:uuid", filesin.DownloadFileService) //下载接口
		file.GET("/List/:pid", filesin.ListFile)                 //文件列表
		file.GET("/Image/:uuid", filesin.TempImages)             //缩略图
		file.POST("/Mkdir/folder", filesin.NewFolder)            //新建文件夹
		file.POST("/RenameFile", filesin.ReNameFile)             //重命名文件
		file.POST("/RenameFolder", filesin.ReNameFolder)         //重命名文件夹
		file.POST("/DelFile", filesin.DelFile)                   //删除文件
		file.POST("/DelFolder", filesin.DelFolder)               //删除文件夹
	}

	if !webConfig.TLS {
		err = r.Run(webConfig.Addr)
		if err != nil {
			return err
		}
	} else {
		err = r.RunTLS(webConfig.Addr, webConfig.TLSCrt, webConfig.TLSKey)
		if err != nil {
			return err
		}
	}
	return nil
}
