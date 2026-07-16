package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := InitDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	r := gin.Default()

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/signup", SignupHandler)
			auth.POST("/login", LoginHandler)
			auth.POST("/logout", AuthMiddleware(), LogoutHandler)
		}

		protected := api.Group("")
		protected.Use(AuthMiddleware())
		{
			protected.GET("/nodes", GetNodesHandler)
			protected.GET("/nodes/:id", GetNodeHandler)
			protected.GET("/map", GetMapHandler)
		}
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
