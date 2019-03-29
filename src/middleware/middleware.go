package middleware

import (
	"strconv"

	"github.com/garyburd/redigo/redis"

	"../config"
	"../constant"
	"../utils"
	"github.com/gin-gonic/gin"
)

// AccessTokenMiddleware : middleware check token
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

	r := config.GetRedisPool().Get()
	defer r.Close()

	uid := utils.ExtractAgentID(authtoken)
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
