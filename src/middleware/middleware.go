package middleware

import (
	"fmt"
	"strconv"

	"github.com/garyburd/redigo/redis"

	. "../config"
	constant "../constant"
	utils "../utils"
	"github.com/gin-gonic/gin"
)

func DummyMiddleware(c *gin.Context) {
	fmt.Println("dummy middle")
	c.Next()
}

func AccessTokenMiddleware(c *gin.Context) {
	authtoken := c.Request.Header.Get("Authorization")

	if authtoken == "" {
		c.AbortWithStatusJSON(400, gin.H{
			"status":       "failure",
			"errorCode":    1000106,
			"errorMessage": "miss auth token",
		})
		return
	}

	r := GetRedisPool().Get()
	defer r.Close()

	uid := utils.ExtractAgentId(authtoken)
	if uid <= 0 {
		c.AbortWithStatusJSON(400, gin.H{
			"status":       "failure",
			"errorCode":    10001048,
			"errorMessage": "invalid token",
		})
		return
	}

	storedToken, err := redis.String(r.Do("GET", strconv.Itoa(int(uid))))
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10001045,
			"errorMessage": "internal server error",
		})
		return
	}

	if storedToken != authtoken {
		c.AbortWithStatusJSON(400, gin.H{
			"status":       "failure",
			"errorCode":    10001044,
			"errorMessage": "Authorization fail",
		})
		return
	}

	// pass token verify
	r.Do("EXPIRE", strconv.Itoa(int(uid)), constant.TokenExpireTime)

	c.Next()
}
