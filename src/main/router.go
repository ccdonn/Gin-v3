package main

import (
	"fmt"

	// . "../config"
	"../middleware"
	"../service"

	"github.com/gin-gonic/gin"
)

func routers() {

	/* testing area */
	api := router.Group("/api")
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
	auth := router.Group("")
	{
		auth.POST("/login", service.Login)
		auth.POST("/resetPassword")
	}

	/* suggestion routing */
	suggestion := router.Group("/suggestion")
	suggestion.Use(middleware.AccessTokenMiddleware)
	{
		suggestion.GET("", service.FindSuggestion)
		suggestion.GET("/:ID", service.GetSuggestion)
		suggestion.POST("", service.CreateSuggestion)
		suggestion.PUT("/:ID", service.PartialUpdateSuggestion)
		// suggestion.DELETE("/:ID", service.DeleteSuggestion)
	}

	/* tutorial routing */
	tutorial := router.Group("/tutorial")
	tutorial.Use(middleware.AccessTokenMiddleware)
	{
		tutorial.GET("", service.FindTutorial)
		tutorial.GET("/:ID", service.GetTutorial)
		tutorial.POST("", service.CreateTutorial)
		tutorial.PUT("/:ID", service.UpdateTutorial)
		tutorial.DELETE("/:ID", service.DeleteTutorial)
	}

	/* brand routing */
	brand := router.Group("/brand")
	brand.Use(middleware.AccessTokenMiddleware)
	{
		brand.GET("")
		brand.PUT("/status/:ID")
		brand.PUT("/hot/:ID")
	}

	/* function routing */
	function := router.Group("/function")
	function.Use(middleware.AccessTokenMiddleware)
	{
		function.GET("/status")
		function.GET("/payment")
		function.PUT("/status")
		function.PUT("/payment")
	}

	/* tip routing */
	tip := router.Group("/tip")
	tip.Use(middleware.AccessTokenMiddleware)
	{
		tip.GET("")
		tip.GET("/:ID")
		tip.PUT("/:ID")
	}

	/* notice routing */
	notice := router.Group("/notice")
	notice.Use(middleware.AccessTokenMiddleware)
	{
		notice.GET("")
		// notice.GET("/:ID")
		notice.POST("")
		// notice.PUT("/:ID")
		notice.DELETE("/:ID")
		notice.PUT("/:ID") // mark as read
	}

	/* push routing */
	push := router.Group("/push")
	push.Use(middleware.AccessTokenMiddleware)
	{
		push.GET("")
		push.POST("")
	}

	/* external service */
	wechat := router.Group("/wechat")
	push.Use(middleware.AccessTokenMiddleware)
	{
		wechat.GET("")
		wechat.PUT("")
		wechat.DELETE("")

		// sub-group
		wechatBrand := wechat.Group("/brand")
		{
			wechatBrand.GET("")
			wechatBrand.POST("")
			wechatBrand.PUT("/:ID")
			wechatBrand.DELETE("/:ID")
		}
	}

	/* external service */
	channel := router.Group("channel")
	channel.Use(middleware.AccessTokenMiddleware)
	{
		channel.GET("/player")
		channel.GET("/agent")
	}
}
