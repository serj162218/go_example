package main

import (
	"fmt"
	"net/http"
	"websocket_with_redis/initializers"
	"websocket_with_redis/websocket"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.ConnectToRedis()
}

func index(c *gin.Context) {
	from := c.Param("from")
	to := c.Param("to")
	fmt.Println("From:" + from + " To:" + to)
	c.HTML(http.StatusOK, "index.html", gin.H{"From": from, "To": to})
}

func main() {
	server := gin.Default()
	server.Static("/assets", "./assets")
	server.LoadHTMLGlob("views/*")

	server.GET("/index/from/:from/to/:to", index)

	ws := websocket.Initialize()
	server.GET("/ws", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		ws.HandleRequest(w, r)
	}))
	server.Run(":8888")
}
