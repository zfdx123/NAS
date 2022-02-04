package ws

import (
	"NAS/src/sysproc"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

func Ws(c *gin.Context) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	//defer conn.Close()
	go wsWrite(conn)
	go wsRead(conn)
}

func wsRead(conn *websocket.Conn) {
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func wsWrite(conn *websocket.Conn) {
	for {
		ws := sysproc.GetWs()
		marshal, _ := json.Marshal(&ws)
		err := conn.WriteMessage(1, marshal)
		if err != nil {
			//log.Println("WS write", err.Error())
			break
		}
		time.Sleep(time.Second * 5)
	}
}
