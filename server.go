package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func upload(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get file err : %s", err.Error()))
		return
	}
	filename := strconv.FormatInt(time.Now().Unix(), 10)
	out, err := os.Create("imgs/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	abs, _ := filepath.Abs("imgs/" + filename)
	res := excute(abs)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "result": res})
}

func main() {
	r := gin.Default()
	r.POST("/upload", upload)
	r.Run(":80")
}
