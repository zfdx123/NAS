package route

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type mess struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func SetUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("系统发生异常!!")
				c.JSON(http.StatusInternalServerError, mess{
					Status:  1,
					Message: "服务器出现异常",
				})
			}
		}()
		c.Next()
	}
}
