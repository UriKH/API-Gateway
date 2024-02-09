package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use()
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
