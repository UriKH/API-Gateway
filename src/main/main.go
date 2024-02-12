package main

import (
	"context"
	"net/http"

	"log"

	"github.com/gin-gonic/gin"

	auth_MS "github.com/TekClinic/Auth-MicroService/generated_files"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func fetch_user_data(c *gin.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth_header := c.GetHeader("Authorization")
		auth_token := auth_header[len("Bearer "):]

		if auth_token == "" {
			c.String(http.StatusInternalServerError, "try again later")
			c.Abort()
		}

		// Authenticate auth_token with auth-microservice
		conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
			c.String(http.StatusInternalServerError, "try again later")
			c.Abort()
		}

		defer conn.Close()
		client := auth_MS.NewAuthServiceClient(conn)

		client_response, err := client.ValidateToken(context.Background(), &auth_MS.TokenRequest{Token: auth_token})
		if err != nil {
			//return error
		}
		user_id := client_response.Id

		//TODO: request user data using user_id from the user data fetcher microservice

		c.String(http.StatusOK, "Hello user_id: %s, your token is %s", user_id, auth_token)
	}
}

func main() {
	router := gin.Default()
	router.Use()

	router.GET("/user/data", fetch_user_data(&gin.Context{}))

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
