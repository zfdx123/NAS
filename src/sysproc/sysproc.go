package sysproc

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
)

var (
	systemctl *exec.Cmd
	ok        = "执行成功"
	no        = "执行失败"
)

type jsonOK struct {
	Status  int    `json:"Status"`
	Message string `json:"Message"`
}
type jsonShell struct {
	Shell string `json:"Shell"`
}

func Shutdown(c *gin.Context) {
	var jShell jsonShell
	err := c.BindJSON(&jShell)
	if err != err {
		c.JSON(http.StatusOK, jsonOK{
			Status:  1,
			Message: "api信息错误",
		})
		return
	}
	switch jShell.Shell {
	case "shutdown":
		if system(jShell.Shell) {
			c.JSON(http.StatusOK, jsonOK{
				Status:  0,
				Message: ok,
			})
			return
		} else {
			c.JSON(http.StatusOK, jsonOK{
				Status:  1,
				Message: no,
			})
			return
		}
	case "reboot":
		if system(jShell.Shell) {
			c.JSON(http.StatusOK, jsonOK{
				Status:  0,
				Message: ok,
			})
			return
		} else {
			c.JSON(http.StatusOK, jsonOK{
				Status:  1,
				Message: ok,
			})
			return
		}
	}
}

func system(shell string) bool {
	arg := []string{shell}
	systemctl = exec.Command("/usr/bin/systemctl", arg...)
	err := systemctl.Run()
	if err != nil {
		return false
	}
	return true
}
