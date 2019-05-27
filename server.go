package main

import (
	"context"
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

var q int

func timeoutMiddleware(timeout time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {

		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)

		defer func() {
			// check if context timeout was reached
			if ctx.Err() == context.DeadlineExceeded {

				// write response and abort the request
				c.Writer.WriteHeader(http.StatusGatewayTimeout)
				c.Abort()
			}

			//cancel to clear resources after finished
			cancel()
		}()

		// replace request with context wrapped request
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

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
func benchMark(c *gin.Context) {
	q++
	c.String(http.StatusOK, "ok, this is No.%d", q)
}

// func drawPic() {
// 	r := image.Rect(0, 0, 1000, 1000)
// 	img := image.NewRGBA(r)
// 	for i := 0; i < 1000; i++ {
// 		for j := 0; j < 1000; j++ {
// 			if i&1 == 1 && j&1 == 1 {
// 				img.Set(i, j, color.RGBA{255, 255, 255, 255})
// 			} else {
// 				img.Set(i, j, color.RGBA{0, 0, 0, 255})
// 			}
// 		}
// 	}

// 	file, _ := os.Create("t.png")
// 	defer file.Close()
// 	png.Encode(file, img)

// }

func main() {
	r := gin.Default()
	r.Use(timeoutMiddleware(60 * time.Second))
	r.POST("/captchsolver", upload)
	r.GET("/benchmark", benchMark)
	r.Run(":80")
}
