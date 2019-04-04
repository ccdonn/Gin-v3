package config

import (
	"../constant"
	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

func GetRedisPool() *redis.Pool {
	if pool == nil {
		pool = newPool()
	}
	return pool
}

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", constant.RedisAddress+":"+constant.RedisPort)
			if err != nil {
				panic(err.Error())
			}
			// password
			if _, err := c.Do("AUTH", constant.RedisPassword); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
	}

}
