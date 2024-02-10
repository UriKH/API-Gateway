package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use()

	router.GET("/user/data", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		authToken := authHeader[len("Bearer "):]
		c.String(http.StatusOK, "Hello %s", authToken)
	})

	router.GET("/users/me", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"data": gin.H{
				"created_at": "2013-12-14T04:35:55.000Z",
				"username":   "TekClinicDev",
				"id":         "2244994945",
				"first_name": "TekClinic",
				"last_name":  "Dev",
			},
		})
	})
	router.Run() // listen and serve on 0.0.0.0:8080
}
