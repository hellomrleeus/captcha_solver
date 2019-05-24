package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get file err : %s", err.Error()))
		return
	}
	filename := header.Filename
	out, err := os.Create("imgs/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	c.String(http.StatusOK, "", "ok")
}

func main() {
	r := gin.Default()
	r.POST("/upload", upload)
	r.Run(":80")
}
