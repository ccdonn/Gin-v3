package middleware

import (
	"log"
	"strconv"

	"github.com/garyburd/redigo/redis"

	"../config"
	"../constant"
	ApiErr "../error"
	"../utils"
	"github.com/gin-gonic/gin"
)

// AccessTokenMiddleware : middleware check token
func AccessTokenMiddleware(c *gin.Context) {
	authtoken := c.Request.Header.Get("Authorization")

	if authtoken == "" {
		c.AbortWithStatusJSON(400, ApiErr.ErrNoToken)
		return
	}

	r := config.GetRedisPool().Get()
	defer r.Close()

	uid := utils.ExtractAgentID(authtoken)
	if uid <= 0 {
		c.AbortWithStatusJSON(400, ApiErr.ErrInvalidToken)
		return
	}

	storedToken, err := redis.String(r.Do("GET", strconv.Itoa(int(uid))))
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, ApiErr.ErrTokenExpire)
		return
	}

	if storedToken != authtoken {
		c.AbortWithStatusJSON(400, ApiErr.ErrAuthFail)
		return
	}

	// pass token verify
	r.Do("EXPIRE", strconv.Itoa(int(uid)), constant.TokenExpireTime)

	c.Next()
}
