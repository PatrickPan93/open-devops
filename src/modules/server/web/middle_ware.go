package web

import "github.com/gin-gonic/gin"

// ConfigMiddleware 设置中间件设置kv值在请求上下文中传递
func ConfigMiddleware(m map[string]interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		for k, v := range m {
			c.Set(k, v)
			c.Next()
		}
	}
}
