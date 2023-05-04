package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
)

func UploadImg(c *gin.Context) {
	filename := uploadSingleImg(c, "upimg")
	c.JSON(200, bson.M{"status": 200, "dir": filename})
}

/*
* the main upload function.
* Save the image with uuid and time-format with (YYYY-MM-DD HH:i:s)
 */
func uploadSingleImg(c *gin.Context, s string) string {
	file, err := c.FormFile(s)
	if err != nil {
		panic(err)
	}
	now := time.Now()
	UUID := uuid.NewString()
	filename := fmt.Sprintf("img_%s_%s.jpg", UUID, now.Format("2006-01-02_150405"))
	c.SaveUploadedFile(file, "./upload/"+filename)
	return filename
}
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
func main() {
	server := gin.Default()
	server.Static("/assets", "./assets")
	server.LoadHTMLGlob("./views/*")

	server.GET("/", Index)
	server.POST("/uploadImg", UploadImg)

	server.Run(":8888")
}
