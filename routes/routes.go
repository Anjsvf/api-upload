package routes

import (
	"github.com/Anjsvf/api-upload/handlers"
	"github.com/Anjsvf/api-upload/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()

	
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, 
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: false,
	}))

	
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	
	v1 := r.Group("/api/v1")
	{
		
		v1.GET("/tags", handlers.GetTags)

	
		v1.GET("/feed", handlers.GetFeed)

		
		posts := v1.Group("/posts")
		posts.Use(middleware.RateLimitByIP())
		{
			posts.POST("", handlers.CreatePost)
		}
	}

	return r
}