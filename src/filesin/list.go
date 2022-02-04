package filesin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type (
	refFiles struct {
		Name string `json:"name"`
		UUID string `json:"UUID"`
		Size int64  `json:"size"`
		PID  string `json:"PID"`
	}
	refFolders struct {
		Name string `json:"name"`
		UUID string `json:"UUID"`
		PID  string `json:"PID"`
	}
)

func ListFile(c *gin.Context) {
	uid := c.GetInt64("uid")
	pid := c.Param("pid")
	tom, f := listSql(uid, pid)
	c.JSON(http.StatusOK, list{
		File:   tom,
		Folder: f,
	})
}
