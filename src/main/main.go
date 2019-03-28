package main

import (
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	router.Run()
}

func init() {
	router = gin.Default()

	// router
	routers()

}
