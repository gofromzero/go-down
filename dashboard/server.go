package dashboard

import "github.com/gin-gonic/gin"

func Server(addr string) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(addr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
