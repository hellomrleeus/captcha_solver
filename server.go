package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func excute(path string) string {

	cmd := exec.Command("/bin/sh", "t.sh", path)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	//log.Println(string(opBytes))
	reg, _ := regexp.Compile(`result:.+`)
	b := reg.Find(opBytes)
	if len(b) < 7 {
		return ""
	}
	b = b[7:]
	return string(b)
}
func upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get file err : %s", err.Error()))
		return
	}
	filename := strconv.FormatInt(time.Now().Unix(), 10) + header.Filename
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
