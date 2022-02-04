package filesin

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func NewFolder(c *gin.Context) {
	uid := c.GetInt64("uid")
	var fol newFolder
	err := c.BindJSON(&fol)

	if err != nil {
		log.Println("解析数据失败", err.Error())
		c.JSON(http.StatusOK, mess{
			Status:  -1,
			Message: "数据解析失败",
		})
		return
	}
	folder, bol := getFolder(uid, fol.PID, fol.Name)
	if bol {
		c.JSON(http.StatusOK, mess{
			Status:  0,
			Message: folder,
		})
	} else {
		c.JSON(http.StatusOK, mess{
			Status:  -1,
			Message: folder,
		})
	}
}

func ReNameFolder(c *gin.Context) {
	var fol reName
	uid := c.GetInt64("uid")
	err := c.BindJSON(&fol)
	if err != nil {
		c.JSON(http.StatusOK, mess{
			Status:  -1,
			Message: "数据解析错误",
		})
		return
	}
	rename, bol := getFolderRename(uid, fol.Name, fol.PID, fol.New)
	if bol {
		c.JSON(http.StatusOK, mess{
			Status:  0,
			Message: rename,
		})
	} else {
		c.JSON(http.StatusOK, mess{
			Status:  -1,
			Message: rename,
		})
	}
}

func DelFolder(c *gin.Context) {
	var fol delFl
	uid := c.GetInt64("uid")
	err := c.BindJSON(&fol)
	if err != nil {
		c.JSON(http.StatusOK, mess{
			Status:  -1,
			Message: "数据解析错误",
		})
		return
	}
	rename, bol := getDelFolder(uid, fol.Name, fol.Pid)
	if bol {
		c.JSON(http.StatusOK, mess{
			Status:  0,
			Message: rename,
		})
	} else {
		c.JSON(http.StatusOK, mess{
			Status:  -1,
			Message: rename,
		})
	}
}
