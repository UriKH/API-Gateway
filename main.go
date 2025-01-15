package main

import (
	"time"

	ginzap "github.com/gin-contrib/zap"

	"go.uber.org/zap"

	ms "github.com/TekClinic/MicroService-Lib"
	"github.com/gin-contrib/location"

	"github.com/TekClinic/API-Gateway/middlewares"
	"github.com/TekClinic/API-Gateway/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	envURIScheme = "URI_SCHEME"
	envURIHost   = "URI_HOST"

	defaultURIScheme = "http"
	defaultURIHost   = "localhost"
	preflightMaxAge  = 12 * time.Hour
)

func main() {
	if ms.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// enable logging
	router.Use(ginzap.Ginzap(zap.L(), time.RFC3339, true))
	// recover in case of panic
	router.Use(gin.Recovery())
	// setup CORS middleware
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Authorization", "Origin", "Content-Length", "Content-Type"},
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		MaxAge:          preflightMaxAge,
	}))
	// setup middleware to discover hostname
	router.Use(location.New(location.Config{
		Scheme: ms.GetOptionalEnv(envURIScheme, defaultURIScheme),
		Host:   ms.GetOptionalEnv(envURIHost, defaultURIHost),
	}))
	// require authorization on all endpoints
	router.Use(middlewares.AuthRequired())

	routes.RegisterPatientRoutes(router)
	routes.RegisterDoctorRoutes(router)
	routes.RegisterAppointmentRoutes(router)
	routes.RegisterTasksRoutes(router)

	err := router.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		zap.L().Fatal("Failed to start server", zap.Error(err))
	}
}
