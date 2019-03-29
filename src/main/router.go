package main

import (
	"fmt"

	// . "../config"
	middleware "../middleware"
	service "../service"

	"github.com/gin-gonic/gin"
)

func routers() {

	api := router.Group("/api")
	api.Use(middleware.DummyMiddleware)
	api.Use(middleware.AccessTokenMiddleware)
	{
		api.GET("/", func(c *gin.Context) {
			fmt.Printf("param is %s\n", c.Query("q"))
			c.JSON(200, gin.H{
				"status": true,
			})
		})
	}

	/* operation without token, login, reset password */
	token := router.Group("/login")
	{
		token.POST("/", service.Login)
	}

	suggestion := router.Group("/suggestion")
	suggestion.Use(middleware.AccessTokenMiddleware)
	{
		suggestion.GET("", service.FindSuggestion)
		suggestion.GET("/:ID", service.GetSuggestion)
		suggestion.POST("", service.CreateSuggestion)
		suggestion.PUT("/:ID", service.PartialUpdateSuggestion)
		// suggestion.DELETE("/:ID", service.DeleteSuggestion)
	}

	tutorial := router.Group("/tutorial")
	tutorial.Use(middleware.AccessTokenMiddleware)
	{
		tutorial.GET("", service.FindTutorial)
		tutorial.GET("/:ID", service.GetTutorial)
		// tutorial.POST("/", service.CreateTutorial)
		// tutorial.PUT("/:ID", service.UpdateTutorial)
		// tutorial.DELETE("/:ID", service.DeleteTutorial)
	}

}
