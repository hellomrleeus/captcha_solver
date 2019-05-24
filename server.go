package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/nfnt/resize"

	"github.com/gin-gonic/gin"
)

func excute(path string) string {

	cmd := exec.Command("/bin/sh", "t.sh", path)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println(err)
	}
	//log.Println(string(opBytes))
	reg, _ := regexp.Compile(`result:.+`)
	b := reg.Find(opBytes)
	if len(b) < 7 {
		return "too small??"
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
	filename := strconv.FormatInt(time.Now().Unix(), 10) + header.Filename + ".png"
	out, err := os.Create("imgs/" + filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Println(err)
		return
	}
	abs, _ := filepath.Abs("imgs/" + filename)
	//adjust height width
	fi, _ := os.Open(abs)

	img, _, err := image.Decode(fi)
	if err != nil {
		fmt.Println(err)
		return
	}

	rz := resize.Resize(250, 70, img, resize.Lanczos3)

	out1, err := os.Create("pngs/resize_" + filename)
	defer out1.Close()
	err = png.Encode(out1, rz)
	if err != nil {
		fmt.Println(err)
		return
	}
	abs, _ = filepath.Abs("pngs/resize_" + filename)
	res := excute(abs)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "result": res})
}

func main() {
	r := gin.Default()
	r.POST("/captchsolver", upload)
	r.Run(":8000")
}
