package main

import (
	"container/list"
	"database/sql"
	"fmt"
	"log"
	"time"

	// . "../config"
	middleware "../middleware"
	. "../response"
	service "../service"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func routers() {

	api := router.Group("/api")
	api.Use(middleware.DummyMiddleware)
	api.Use(middleware.AccessTokenMiddleware)
	{
		api.GET("/", func(c *gin.Context) {
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
		suggestion.GET("/", service.FindSuggestion)
		suggestion.GET("/:ID", service.GetSuggestion)
		suggestion.POST("/", service.CreateSuggestion)
		suggestion.PUT("/:ID", service.PartialUpdateSuggestion)
		// suggestion.DELETE("/:ID", service.DeleteSuggestion)
	}

	tutorial := router.Group("/tutorial/:ID")
	{
		tutorial.POST("/", func(c *gin.Context) {

			c.JSON(200, gin.H{
				"message": "message222",
				"id":      c.Param("ID"),
			})
		})
	}

	me := router.Group("/me")
	{
		me.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"name": "name",
			})
		})

		me.GET("/dbdirect", func(c *gin.Context) {
			db, err := sql.Open("mysql", "root:root@/promotion")
			log.Println(db)
			rows, err := db.Query("select agent_id, nickname, username from account where agent_id < ?", 1000)
			if err != nil {
				log.Fatal("sql error")
			}
			fmt.Println(rows)
			var (
				agentID  int32
				nickname string
				username string
			)

			res := list.New()
			for rows.Next() {
				err := rows.Scan(&agentID, &nickname, &username)
				if err != nil {
					log.Fatal("row scan error")
					c.JSON(500, "Error")
					return
				}

				res.PushBack(MeResponse{
					AgentID:  agentID,
					Username: username,
					Nickname: nickname,
				})

			}

			slice := make([]MeResponse, res.Len())

			cnt := 0
			for iter := res.Front(); iter != nil; iter = iter.Next() {
				slice[cnt] = iter.Value.(MeResponse)
				cnt++
			}

			c.JSON(200, gin.H{
				"status":    "success",
				"Timestamp": time.Now(),
				// "result": MeResponse{
				// 	AgentID:  agentID,
				// 	Username: username,
				// 	Nickname: nickname,
				// },
				"size":   cnt,
				"result": slice,
			})

		})

		me.GET("/db", func(c *gin.Context) {
			// db := GetConn()
			db, err := sql.Open("mysql", "root:root@/promotion")
			if err != nil {
				c.JSON(500, gin.H{
					"error": "db connection error",
				})
				return
			}

			rows, err := db.Query("select agent_id, nickname, username from account where agent_id = ?", 19)
			if err != nil {
				// log.Panic(err)
				c.JSON(500, gin.H{
					"error": "sql error",
				})
				return
			}

			var (
				agentID  int32
				nickname string
				username string
			)

			log.Println(rows)

			if rows.Next() {
				err = rows.Scan(&agentID, &nickname, &username)
				if err != nil {
					// log.Panic("scan error")
					c.JSON(500, gin.H{
						"error": "scan error",
					})
					return
				}
			} else {
				c.JSON(200, gin.H{
					"status":    "success",
					"Timestamp": time.Now(),
					"size":      0,
					"result":    nil,
				})
				return
			}

			c.JSON(200, gin.H{
				"status":    "success",
				"Timestamp": time.Now(),
				"size":      1,
				"result": MeResponse{
					AgentID:  agentID,
					Username: username,
					Nickname: nickname,
				},
			})
		})
	}

}
