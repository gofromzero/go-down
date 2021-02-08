package dashboard

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Routers() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())
	r.Use(CorsMiddleware())
	r.Use(HandleOrigin())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": 0,
			"msg":    "ok",

			"data": gin.H{
				"count": 0,
				"rows":  []int{},
			},
		})
	})
	return r
}

// cors
func CorsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:        true,
		AllowMethods:           []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:           []string{"Content-Type", "AccessToken", "X-CSRF-Token", "Authorization", "Token", "X-Token", "user-id"},
		AllowCredentials:       false,
		ExposeHeaders:          []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		MaxAge:                 12 * time.Hour,
		AllowWildcard:          false,
		AllowBrowserExtensions: false,
		AllowWebSockets:        false,
	})
}

// origin
// 处理跨域请求,支持options访问
func HandleOrigin() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		//放行所有OPTIONS方法
		if method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		// 处理请求
		c.Next()
	}
}
