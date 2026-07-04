package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/raaj2493/Shortly/Backend/internal/handlers"
)

func SetupRouters (r *gin.Engine , authHandler *handlers.AuthHandler){
	// Create our version 1 group base path
	v1 := r.Group("/api/v1")
	{
		// Authentication Sub-Group
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.RegisterUser)
			// Future routes will drop in here cleanly:
			// auth.POST("/login", authHandler.LoginUser)
		}

		// Future URL Shortener Sub-Group will live here contextually:
		// urls := v1.Group("/urls")
		// {
		//     urls.POST("/shorten", urlHandler.ShortenURL)
		// }
	}
}