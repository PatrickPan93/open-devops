package web

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func StartGin(httpAddr string) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	// 设置路由
	configRoutes(r)
	s := &http.Server{
		Addr:           httpAddr,
		Handler:        r,
		ReadTimeout:    time.Second * 5,
		WriteTimeout:   time.Second * 5,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("web server available: %s\n", httpAddr)
	return s.ListenAndServe()

}
