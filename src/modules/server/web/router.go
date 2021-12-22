package web

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func configRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})
		api.GET("now-ts", GetNowTs)
		api.POST("/node-path", NodePathAdd)
		api.GET("/node-path", NodePathQuery)
		api.DELETE("/node-path", NodePathDelete)
		api.POST("/resource-mount", ResourceMount)
	}
}

func GetNowTs(c *gin.Context) {
	c.String(http.StatusOK, time.Now().Format("2006-01-02 15:04:05"))
}
