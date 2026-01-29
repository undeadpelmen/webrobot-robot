package http

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/undeadpelmen/webrobot-robot/hardware/robots"
)

var (
	outch chan string
	errch chan error

	status *string
)

func RobotHttpFunc(stat *string, out chan string, errc chan error) {
	outch = out
	errch = errc

	status = stat

	gin.DisableConsoleColor()
	gin.DisableBindValidation()

	f, _ := os.Create(filepath.Join(os.TempDir(), "webrobot/gin.log"))

	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	r := gin.Default()

	r.Use(CORSMiddleware())

	r.GET("/ping", pong)
	r.POST("/robot/api/command", commandHandler)
	r.GET("/robot/api/status", statusHandler)

	err := r.Run(":8080")
	if err != nil {
		errch <- err
	}

}

func pong(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func commandHandler(c *gin.Context) {
	var b body
	if err := c.BindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !robots.ValidCommand(b.Command) {
		fmt.Println(b.Command)
		fmt.Println(b)
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid command", "command": b.Command})
		return
	}

	outch <- b.Command

	c.JSON(http.StatusOK, gin.H{"data": b})
}

func statusHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": *status})
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
